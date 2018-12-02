#define _GENFUN(x, y) x##y
#define GENFUN(x, y) _GENFUN(x, y)

void GENFUN(multiplyConjugateInline, __SUBARCH__)(float *vecA, float *vecB, unsigned int length) {

    float *aPtr = (float *)vecA;
    float *bPtr = (float *)vecB;

    for (unsigned int i = 0; i < length; i++) {
      // A = (Ar) + (Ai)i
      // B = (Br) + (Bi)i
      // ( Ar * Br − Ai * Bi) + ( Ar * Bi + Ai * Br) i => Complex Multiplication

      float Ar = vecA[i*2+0];
      float Ai = vecA[i*2+1];
      float Br = vecB[i*2+0];
      float Bi = vecB[i*2+1];

      Bi = -Bi; // Conjugated

      vecA[i*2+0] = (Ar * Br) - (Ai * Bi);
      vecA[i*2+1] = (Ar * Bi) + (Ai * Br);
    }
}