#define _GENFUN(x, y) x##y
#define GENFUN(x, y) _GENFUN(x, y)

void GENFUN(dotProductComplexComplex, __SUBARCH__)(float *result, float *input, const float *taps, unsigned int length) {
    unsigned int cBlocks = length / 16;
    float sum0[2] = {0, 0};
    float sum1[2] = {0, 0};

    float *iPtr = (float *)input;
    float *tPtr = (float *)taps;

    for (unsigned int i = 0; i < cBlocks; i++) {
        sum0[0] += iPtr[0] * tPtr[0] - iPtr[1] * tPtr[1];
        sum0[1] += iPtr[0] * tPtr[1] + iPtr[1] * tPtr[0];
        sum1[0] += iPtr[2] * tPtr[2] - iPtr[3] * tPtr[3];
        sum1[1] += iPtr[2] * tPtr[3] + iPtr[3] * tPtr[2];
        iPtr += 4;
        tPtr += 4;
    }

    result[0] = sum0[0] + sum1[0];
    result[1] = sum0[1] + sum1[1];
}