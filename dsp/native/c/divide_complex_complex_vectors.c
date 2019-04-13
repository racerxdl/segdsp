#define _GENFUN(x, y) x##y
#define GENFUN(x, y) _GENFUN(x, y)

void GENFUN(divideComplexComplexVectors, __SUBARCH__)(float *A, float *B, unsigned int length) {
    for (unsigned int i = 0; i < length; i++) {
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

      float Ar = A[i*2+0];
      float Ai = A[i*2+1];
      float Br = B[i*2+0];
      float Bi = B[i*2+1];

      float denom = Br * Br + Bi * Bi;

      A[i*2+0] = (Ar * Br + Ai * Bi) / denom;
      A[i*2+1] = (Br * Ai - Ar * Bi) / denom;
    }
}



