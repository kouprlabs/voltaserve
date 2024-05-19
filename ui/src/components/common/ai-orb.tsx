/** @jsxImportSource @emotion/react */
import { css, keyframes } from '@emotion/react'

const rotate = keyframes`
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
`

const animate = keyframes`
  0% { transform: scale(1) rotate(0deg); }
  50% { transform: scale(1.2) rotate(180deg); }
  100% { transform: scale(1) rotate(360deg); }
`

const orbStyle = css`
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: radial-gradient(
    circle,
    rgba(255, 255, 255, 1) 0%,
    rgba(0, 0, 0, 1) 100%
  );
  position: relative;
  overflow: hidden;
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

const AiOrb = () => {
  return <div css={orbStyle}></div>
}

export default AiOrb
