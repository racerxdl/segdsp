#define _GENFUN(x, y) x##y
#define GENFUN(x, y) _GENFUN(x, y)

#include <math.h>

inline void multiplyComplex(float *output, float *a, float *b) {
      float Ar = a[0];
      float Ai = a[1];
      float Br = b[0];
      float Bi = b[1];

      output[0] = (Ar * Br) - (Ai * Bi);
      output[1] = (Ar * Bi) + (Ai * Br);
}

#if __amd64__ || __x86_64__ || __X86__ || __i386__
// Use hardware SQRT
inline volatile double asmSqrt(double a){
    double b = 0;
    asm volatile(
        "movq %1, %%xmm0 \n"
        "sqrtsd %%xmm0, %%xmm1 \n"
        "movq %%xmm1, %0 \n"
        : "=r"(b)
        : "g"(a)
        : "xmm0", "xmm1", "memory"
    );
    return b;
}
#else
// Fallback
inline volatile double asmSqrt(double a){
    return sqrt(a);
}
#endif

inline void normalizeComplex(float *complex) {
  float mag = asmSqrt(complex[0] * complex[0] + complex[1] * complex[1]);
  complex[0] /= mag;
  complex[1] /= mag;
}

void GENFUN(rotateComplex, __SUBARCH__)(float *input, float *output, float *phaseIncrement, float *phase, unsigned int length) {
    for (unsigned int i = 0; i < length / 512; i++) {
        for (unsigned int j = 0; j < 512; j++) {
            multiplyComplex(output, input, phase);
            output += 2;
            input += 2;
            multiplyComplex(phase, phaseIncrement, phase);
        }
        normalizeComplex(phase);
    }

    for (unsigned int i = 0; i < length % 512; i++) {
        multiplyComplex(output, input, phase);
        output += 2;
        input += 2;
        multiplyComplex(phase, phase, phaseIncrement);
    }
}
