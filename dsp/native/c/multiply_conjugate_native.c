#define _GENFUN(x, y) x##y
#define GENFUN(x, y) _GENFUN(x, y)

void GENFUN(multiplyConjugate, __SUBARCH__)(float *vecA, float *vecB, float *output, unsigned int length) {
    for (unsigned int i = 0; i < length; i++) {
      // (a+bi) (c+di) = (ac−bd) + (ad+bc)i.
      // A = (Ar) + (Ai)i
      // B = (Br) + (Bi)i
      // ( Ar * Br − Ai * Bi) + ( Ar * Bi + Ai * Br) i => Complex Multiplication

      float Ar = vecA[i*2+0];
      float Ai = vecA[i*2+1];
      float Br = vecB[i*2+0];
      float Bi = vecB[i*2+1];

      Bi = -Bi; // Conjugated

      output[i*2+0] = (Ar * Br) - (Ai * Bi);
      output[i*2+1] = (Ar * Bi) + (Ai * Br);
    }
}
