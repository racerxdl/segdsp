
#if defined __GNUC__
#  define SIMD_ALIGNED(x) __attribute__((aligned(x)))
#elif defined __clang__
#  define SIMD_ALIGNED(x) __attribute__((aligned(x)))
#elif _MSC_VER
#  define SIMD_ALIGNED(x) __declspec(align(x))
#else
#  define SIMD_ALIGNED(x)
#endif
