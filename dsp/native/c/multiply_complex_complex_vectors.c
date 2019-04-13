#define _GENFUN(x, y) x##y
#define GENFUN(x, y) _GENFUN(x, y)

void GENFUN(multiplyComplexComplexVectors, __SUBARCH__)(float *A, float *B, unsigned int length) {
    for (unsigned int i = 0; i < length; i++) {
      // (a+bi) (c+di) = (ac−bd) + (ad+bc)i.
      // A = (Ar) + (Ai)i
      // B = (Br) + (Bi)i
      // ( Ar * Br − Ai * Bi) + ( Ar * Bi + Ai * Br) i => Complex Multiplication

      float Ar = A[i*2+0];
      float Ai = A[i*2+1];
      float Br = B[i*2+0];
      float Bi = B[i*2+1];

      A[i*2+0] = (Ar * Br) - (Ai * Bi);
      A[i*2+1] = (Ar * Bi) + (Ai * Br);
    }
}
