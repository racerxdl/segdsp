#define _GENFUN(x, y) x##y
#define GENFUN(x, y) _GENFUN(x, y)

void GENFUN(subtractComplexComplexVectors, __SUBARCH__)(float *A, float *B, unsigned int length) {
    unsigned int floatNumbers = length * 2; // Each complex contains 2 floats
    unsigned int cBlocks = floatNumbers / 4;
    unsigned int c = 0;

    // Unfold by 4
    for (unsigned int i = 0; i < cBlocks; i++) {
        A[i*4+0] = A[i*4+0] - B[i*4+0];
        A[i*4+1] = A[i*4+1] - B[i*4+1];
        A[i*4+2] = A[i*4+2] - B[i*4+2];
        A[i*4+3] = A[i*4+3] - B[i*4+3];
        c+=4;
    }

    // Add remaning non multiple of 4
    for (int i = c; i < floatNumbers; i++) {
        A[i] = A[i] - B[i];
    }
}
