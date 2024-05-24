/** @jsxImportSource @emotion/react */
import { css, keyframes } from '@emotion/react'

const rotate = keyframes`
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
`

const animate = keyframes`
  0% { transform: scale(1) rotate(0deg); background: linear-gradient(45deg, rgba(255,0,0,0.5), rgba(255,255,0,0.5), rgba(0,255,0,0.5), rgba(0,255,255,0.5), rgba(0,0,255,0.5), rgba(255,0,255,0.5)); }
  25% { transform: scale(1.05) rotate(90deg); background: linear-gradient(45deg, rgba(255,0,0,0.55), rgba(255,255,0,0.55), rgba(0,255,0,0.55), rgba(0,255,255,0.55), rgba(0,0,255,0.55), rgba(255,0,255,0.55)); }
  50% { transform: scale(1) rotate(180deg); background: linear-gradient(45deg, rgba(255,0,0,0.6), rgba(255,255,0,0.6), rgba(0,255,0,0.6), rgba(0,255,255,0.6), rgba(0,0,255,0.6), rgba(255,0,255,0.6)); }
  75% { transform: scale(1.05) rotate(270deg); background: linear-gradient(45deg, rgba(255,0,0,0.55), rgba(255,255,0,0.55), rgba(0,255,0,0.55), rgba(0,255,255,0.55), rgba(0,0,255,0.55), rgba(255,0,255,0.55)); }
  100% { transform: scale(1) rotate(360deg); background: linear-gradient(45deg, rgba(255,0,0,0.5), rgba(255,255,0,0.5), rgba(0,255,0,0.5), rgba(0,255,255,0.5), rgba(0,0,255,0.5), rgba(255,0,255,0.5)); }
`

const getOrbStyle = (width: string, height: string) => css`
  width: ${width};
  height: ${height};
  border-radius: 50%;
  background: radial-gradient(
    circle,
    rgba(255, 255, 255, 1) 0%,
    rgba(0, 0, 0, 1) 100%
  );
  position: relative;
  overflow: hidden;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.5);
  animation: ${rotate} 10s linear infinite;
  &:before,
  &:after {
    content: '';
    position: absolute;
    width: 200%;
    height: 200%;
    top: -50%;
    left: -50%;
    background: linear-gradient(
      45deg,
      rgba(255, 0, 0, 0.5),
      rgba(255, 255, 0, 0.5),
      rgba(0, 255, 0, 0.5),
      rgba(0, 255, 255, 0.5),
      rgba(0, 0, 255, 0.5),
      rgba(255, 0, 255, 0.5)
    );
    animation: ${animate} 6s ease-in-out infinite;
  }
  &:after {
    animation-direction: reverse;
  }
`

export type OrbProps = {
  width: string
  height: string
}

const Orb = ({ width, height }: OrbProps) => {
  return <div css={getOrbStyle(width, height)}></div>
}

export default Orb
