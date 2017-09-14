#include <stdlib.h>
#include <errno.h>
#include <uv.h>
#include <netdb.h>
#include <arpa/inet.h>
#include <sys/stat.h>
#include <fcntl.h>

#ifndef _WIN32
#include <sys/socket.h>
#include <unistd.h>
#endif

#define ASSERT(expr)                                           \
    do                                                         \
    {                                                          \
        if (!(expr))                                           \
        {                                                      \
            fprintf(stderr,                                    \
                    "Assertion failed in %s on line %d: %s\n", \
                    __FILE__,                                  \
                    __LINE__,                                  \
                    #expr);                                    \
            abort();                                           \
        }                                                      \
    } while (0)

#define UV_SIZEOF_SOCKADDR_IN ((int)sizeof(struct sockaddr_in))

extern void __uv_connect_cb(uv_connect_t *req, int status);
extern void __uv_connection_cb(uv_stream_t *stream, int status);
extern void __uv_write_cb(uv_write_t *req, int status);
extern void __uv_read_cb(uv_stream_t *stream, ssize_t nread, uv_buf_t *buf);
extern void __uv_udp_recv_cb(uv_udp_t *handle, ssize_t nread, uv_buf_t *buf, struct sockaddr *addr, unsigned flags);
extern void __uv_udp_send_cb(uv_udp_send_t *req, int status);
extern void __uv_timer_cb(uv_timer_t *timer, int status);
extern void __uv_poll_cb(uv_poll_t *p, int status, int events);
extern void __uv_signal_cb(uv_signal_t *s, int signum);
extern void __uv_idle_cb(uv_idle_t *handle, int status);
extern void __uv_close_cb(uv_handle_t *handle);
extern void __uv_prepare_cb(uv_prepare_t *handle);
extern void __uv_async_cb(uv_prepare_t *handle);
extern void __uv_check_cb(uv_check_t *handle);
extern void __uv_shutdown_cb(uv_shutdown_t *req, int status);
extern void __uv_exit_cb(uv_process_t *process, int exit_status, int term_signal);

typedef struct connection_context_s
{
    uv_poll_t poll_handle;
    uv_timer_t timer_handle;
    uv_os_sock_t sock;
    size_t read, sent;
    int is_server_connection;
    int open_handles;
    int got_fin, sent_fin;
    unsigned int events, delayed_events;
} connection_context_t;

static void _uv_alloc_cb(uv_handle_t *handle, size_t suggested_size, uv_buf_t *buf)
{
    char *base;
    base = (char *)calloc(1, suggested_size);
    if (!base)
        *buf = uv_buf_init(NULL, 0);
    else
        *buf = uv_buf_init(base, suggested_size);
}

static uv_buf_t *uv_buf_malloc(unsigned int len)
{
    uv_buf_t *buf;
    buf = malloc(len * sizeof(uv_buf_t));
    return buf;
}

static void uv_buf_set(uv_buf_t *bufs, uint index, uv_buf_t buf)
{
    bufs[index] = buf;
}

static int _uv_udp_send(uv_udp_send_t *req, uv_udp_t *handle, uv_buf_t bufs[], unsigned int bufcnt, struct sockaddr *addr)
{
    return uv_udp_send(req, handle, bufs, bufcnt, addr, __uv_udp_send_cb);
}

static int _uv_udp_recv_start(uv_udp_t *udp)
{
    return uv_udp_recv_start(udp, _uv_alloc_cb, __uv_udp_recv_cb);
}

static int _uv_tcp_connect(uv_connect_t *req, uv_tcp_t *handle, struct sockaddr *address)
{
    return uv_tcp_connect(req, handle, address, __uv_connect_cb);
}

static void _uv_pipe_connect(uv_connect_t *req, uv_pipe_t *handle, const char *name)
{
    uv_pipe_connect(req, handle, name, __uv_connect_cb);
}

static int _uv_listen(uv_stream_t *stream, int backlog)
{
    return uv_listen(stream, backlog, __uv_connection_cb);
}

static int _uv_read_start(uv_stream_t *stream)
{
    return uv_read_start(stream, _uv_alloc_cb, __uv_read_cb);
}

static int _uv_write(uv_write_t *req, uv_stream_t *handle, uv_buf_t bufs[], int bufcnt)
{
    return uv_write(req, handle, bufs, bufcnt, __uv_write_cb);
}

static int _uv_write2(uv_write_t *req, uv_stream_t *handle, uv_buf_t bufs[], int bufcnt, uv_stream_t *send_handle)
{
    return uv_write2(req, handle, bufs, bufcnt, send_handle, __uv_write_cb);
}

static int _uv_try_write(uv_stream_t *handle, uv_buf_t bufs[], int bufcnt)
{
    return uv_try_write(handle, bufs, bufcnt);
}

static void _uv_close(uv_handle_t *handle)
{
    uv_close(handle, __uv_close_cb);
}

static int _uv_shutdown(uv_shutdown_t *req, uv_stream_t *handle)
{
    return uv_shutdown(req, handle, __uv_shutdown_cb);
}

static int _uv_timer_start(uv_timer_t *timer, uint64_t timeout, uint64_t repeat)
{
    return uv_timer_start(timer, __uv_timer_cb, timeout, repeat);
}

static int _uv_poll_start(uv_poll_t *p, int events)
{
    return uv_poll_start(p, events, __uv_poll_cb);
}

static int _uv_signal_start(uv_signal_t *s, int sigNum)
{
    return uv_signal_start(s, __uv_signal_cb, sigNum);
}

static int _uv_signal_start_oneshot(uv_signal_t *s, int sigNum)
{
    return uv_signal_start_oneshot(s, __uv_signal_cb, sigNum);
}

static int _uv_prepare_start(uv_prepare_t *handle)
{
    return uv_prepare_start(handle, __uv_prepare_cb);
}

static int _uv_async_init(uv_loop_t *loop, uv_async_t *handle)
{
    return uv_async_init(loop, handle, __uv_async_cb);
}

static int _uv_check_start(uv_check_t *handle)
{
    return uv_check_start(handle, __uv_check_cb);
}

static int _uv_idle_start(uv_idle_t *idle)
{
    return uv_idle_start(idle, __uv_idle_cb);
}

static int _uv_spawn(uv_loop_t *loop, uv_process_t *process, uv_process_options_t *options)
{
    options->exit_cb = __uv_exit_cb;
    return uv_spawn(loop, process, options);
}

static uv_os_sock_t create_socket(struct sockaddr_in *bind_addr, int bound_socket, int protocol)
{
    uv_os_sock_t sock;
    int r;

    sock = socket(AF_INET, SOCK_STREAM, protocol);

#ifdef _WIN32
    ASSERT(sock != INVALID_SOCKET);
#else
    ASSERT(sock >= 0);
#endif

#ifndef _WIN32
    {
        /* Allow reuse of the port. */
        int yes = 1;
        r = setsockopt(sock, SOL_SOCKET, SO_REUSEADDR, (char *)&yes, sizeof yes);
        ASSERT(r == 0);
    }
#endif

    if (bound_socket == 1)
    {
        r = bind(sock, (const struct sockaddr *)bind_addr, sizeof *bind_addr);
        ASSERT(r == 0);
    }

    return sock;
}

static uv_os_sock_t create_tcp_socket(struct sockaddr_in *bind_addr, int bound_socket)
{
    uv_os_sock_t sock;
    sock = create_socket(bind_addr, bound_socket, IPPROTO_TCP);

#ifndef _WIN32
    {
        int yes = 1;
        ASSERT(setsockopt(sock, IPPROTO_TCP, TCP_NODELAY, (char *)&yes, sizeof yes) == 0);
    }
#endif

    return sock;
}

static int connect_socket(uv_os_sock_t sock, struct sockaddr *saddr)
{
    return connect(sock, saddr, sizeof(struct sockaddr));
}

static int close_socket(uv_os_sock_t sock)
{
    int r;
#ifdef _WIN32
    r = closesocket(sock);
#else
    r = close(sock);
#endif
    return r;
}

static connection_context_t *create_connection_context(uv_os_sock_t sock, int is_server_connection)
{
    int r;
    connection_context_t *context;

    context = (connection_context_t *)malloc(sizeof *context);
    ASSERT(context != NULL);

    context->sock = sock;
    context->is_server_connection = is_server_connection;
    context->read = 0;
    context->sent = 0;
    context->open_handles = 0;
    context->events = 0;
    context->delayed_events = 0;
    context->got_fin = 0;
    context->sent_fin = 0;

    return context;
}

static void test_sendAndRecv(uv_os_sock_t sockfd)
{
    int numbytes;
    char buf[4096];

    if (send(sockfd, "Hello world from sock client!\n", 14, 0) == -1)
    {
        return;
    }
    printf("After the send function \n");

    if ((numbytes = recv(sockfd, buf, 4096, 0)) != -1)
    {
        buf[numbytes] = '\0';

        printf("Received in pid=%d, text=: %s \n", getpid(), buf);
    }
    else
    {
        printf("recv error\n");
    }
}

static int test_Open(char *path)
{
    return open(path, O_RDWR);
}
