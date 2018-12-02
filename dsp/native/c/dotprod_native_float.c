#define _GENFUN(x, y) x##y
#define GENFUN(x, y) _GENFUN(x, y)

void GENFUN(dotProductFloat, __SUBARCH__)(float *result, float *input, const float *taps, unsigned int length) {
    float res = 0;

    float *iPtr = (float *)input;
    float *tPtr = (float *)taps;

    for (unsigned int i = 0; i < length; i++) {
      res += ((*iPtr++) * (*tPtr++));
    }

    *result = res;
}
