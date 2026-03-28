#include <sycl/sycl.hpp>

struct point {
    double X;
    double Y;
};

extern "C" SYCL_EXT_ONEAPI_FUNCTION_PROPERTY((sycl::ext::oneapi::experimental::nd_range_kernel<1>))
void quadratic(
    double x0, double y0, double x1, double y1,
    double x2, double y2, double ds, struct point p[]
) {
    auto item = sycl::ext::oneapi::this_work_item::get_nd_item<1>();
    auto idx = item.get_global_id(0);
    auto t = (double)idx / ds;
    auto u = 1 - t;
    auto a = u * u;
    auto b = 2 * u * t;
    auto c = t * t;
    p[idx].X = a*x0+b*x1+c*x2;
    p[idx].Y = a*y0+b*y1+c*y2;
}

extern "C" SYCL_EXT_ONEAPI_FUNCTION_PROPERTY((sycl::ext::oneapi::experimental::nd_range_kernel<1>))
void cubic(
    double x0, double y0, double x1, double y1,
    double x2, double y2, double x3, double y3,
    double ds, struct point p[]
) {
    auto item = sycl::ext::oneapi::this_work_item::get_nd_item<1>();
    auto idx = item.get_global_id(0);
    auto t = (double)idx / ds;
    auto u = 1 - t;
    auto a = u * u * u;
    auto b = 3 * u * u * t;
    auto c = 3 * u * t * t;
    auto d = t * t * t;
    p[idx].X = a*x0+b*x1+c*x2+d*x3;
    p[idx].Y = a*y0+b*y1+c*y2+d*y3;
}
