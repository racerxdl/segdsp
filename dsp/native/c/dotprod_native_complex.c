#define _GENFUN(x, y) x##y
#define GENFUN(x, y) _GENFUN(x, y)

void GENFUN(dotProductComplex, __SUBARCH__)(float *result, float *input, const float *taps, unsigned int length) {
    float res[2] = {0, 0};

    float *iPtr = (float *)input;
    float *tPtr = (float *)taps;

    for (unsigned int i = 0; i < length; i++) {
      res[0] += ((*iPtr++) * (*tPtr));
      res[1] += ((*iPtr++) * (*tPtr++));
    }

    result[0] = res[0];
    result[1] = res[1];
}