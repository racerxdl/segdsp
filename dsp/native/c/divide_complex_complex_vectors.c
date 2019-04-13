#define _GENFUN(x, y) x##y
#define GENFUN(x, y) _GENFUN(x, y)

inline void Divide(float *A, float *B) {
    // (a+bi) / (c+di) = ((ac+bd) + (cb-ad)i) / ((c*c) + (d*d))
    // A = (Ar) + (Ai)i
    // B = (Br) + (Bi)i
    // a = Ar
    // b = Ai
    // c = Br
    // d = Bi
    // ((Ar*Br+Ai*Bi) + (Ar*Bi+Br*Ai)i) / ((Br*Br) + (Bi*Bi))
    //
    // denom = ((Br * Br) + (Bi * Bi))
    // ((Ar * Br + Ai * Bi) + (Ar * Bi + Br * Ai)i) / denom => Complex Divison

    float Ar = A[0];
    float Ai = A[1];
    float Br = B[0];
    float Bi = B[1];

    float denom = Br * Br + Bi * Bi;

    A[0] = (Ar * Br + Ai * Bi) / denom;
    A[1] = (Br * Ai - Ar * Bi) / denom;
}

void GENFUN(divideComplexComplexVectors, __SUBARCH__)(float *A, float *B, unsigned int length) {
    const unsigned int cBlocks = length / 4;
    unsigned int c = 0;

    float *aPtr = A;
    float *bPtr = B;

    for (unsigned int i = 0; i < cBlocks; i++) {
      Divide(aPtr+0, bPtr+0);
      Divide(aPtr+2, bPtr+2);
      Divide(aPtr+4, bPtr+4);
      Divide(aPtr+6, bPtr+6);
      aPtr += 8;
      bPtr += 8;
      c += 4;
    }

    for (unsigned int i = c; i < length; i++) {
      Divide(aPtr, bPtr);
      aPtr += 2;
      bPtr += 2;
    }
}



