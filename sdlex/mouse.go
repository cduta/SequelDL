package sdlex

func IsMouseMotionState(stateBitMask, state uint32) bool {
  return (stateBitMask >> (state - 1)) & 1 > 0
}