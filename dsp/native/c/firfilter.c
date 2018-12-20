#include "helper.h"

#define _GENFUN(x, y) x##y
#define GENFUN(x, y) _GENFUN(x, y)

inline void __dotProduct__(float *result, float *input, const float *taps, unsigned int length) {
    float res[4] = {0, 0, 0, 0};

    float *iPtr = (float *)input;
    float *tPtr = (float *)taps;

    for (unsigned int i = 0; i < length / 2; i++) {
      res[0] += ((*iPtr++) * (*tPtr));
      res[1] += ((*iPtr++) * (*tPtr++));
      res[2] += ((*iPtr++) * (*tPtr));
      res[3] += ((*iPtr++) * (*tPtr++));
    }

    result[0] = res[0] + res[2];
    result[1] = res[1] + res[3];
}

void GENFUN(firFilter, __SUBARCH__)(float *result, float *input, const float *taps, unsigned int lengthTaps, unsigned int length) {
  for (unsigned int i = 0; i < length - lengthTaps; i++) {
    __dotProduct__(&result[i], &input[i], taps, lengthTaps);
  }
}

void GENFUN(firFilterDecimate, __SUBARCH__)(float *result, float *input, const float *taps, unsigned int decimate, unsigned int lengthTaps, unsigned int length) {
  unsigned int j = 0;
  for (unsigned int i = 0; i < length - lengthTaps; i++) {
    __dotProduct__(&result[i], &input[j], taps, lengthTaps);
    j += decimate;
  }
}
