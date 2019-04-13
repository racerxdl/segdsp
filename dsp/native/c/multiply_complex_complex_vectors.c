#define _GENFUN(x, y) x##y
#define GENFUN(x, y) _GENFUN(x, y)

inline void Multiply(float *A, float *B) {
    // (a+bi) (c+di) = (ac−bd) + (ad+bc)i.
    // A = (Ar) + (Ai)i
    // B = (Br) + (Bi)i
    // ( Ar * Br − Ai * Bi) + ( Ar * Bi + Ai * Br) i => Complex Multiplication
    float Ar = A[0];
    float Ai = A[1];
    float Br = B[0];
    float Bi = B[1];

    A[0] = (Ar * Br) - (Ai * Bi);
    A[1] = (Ar * Bi) + (Ai * Br);
}

void GENFUN(multiplyComplexComplexVectors, __SUBARCH__)(float *A, float *B, unsigned int length) {
    const unsigned int cBlocks = length / 4;
    unsigned int c = 0;

    float *aPtr = A;
    float *bPtr = B;

    for (unsigned int i = 0; i < cBlocks; i++) {
      Multiply(aPtr+0, bPtr+0);
      Multiply(aPtr+2, bPtr+2);
      Multiply(aPtr+4, bPtr+4);
      Multiply(aPtr+6, bPtr+6);
      aPtr += 8;
      bPtr += 8;
      c += 4;
    }

    for (unsigned int i = c; i < length; i++) {
      Multiply(aPtr, bPtr);
      aPtr += 2;
      bPtr += 2;
    }
}
