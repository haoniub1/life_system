<template>
  <n-config-provider :theme="darkTheme">
    <n-message-provider>
      <n-dialog-provider>
        <!-- 水印背景 -->
        <div class="watermark">
          <div class="watermark-text">上士闻道</div>
          <div class="watermark-text">勤而行之</div>
        </div>
        <!-- 主内容 -->
        <div class="main-content">
          <router-view />
        </div>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { NConfigProvider, NMessageProvider, NDialogProvider, darkTheme } from 'naive-ui'
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  background-color: #1a1a1a;
  color: #e0e0e0;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

#app {
  width: 100%;
  height: 100vh;
  position: relative;
}

/* 水印层 - 固定背景 */
.watermark {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  pointer-events: none;
  z-index: 0;
  user-select: none;
}

.watermark-text {
  font-size: 4rem;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.03);
  letter-spacing: 1rem;
  font-family: 'KaiTi', 'STKaiti', 'SimSun', serif;
  line-height: 1.5;
}

/* 主内容层 - 在水印上方 */
.main-content {
  position: relative;
  z-index: 1;
  width: 100%;
  height: 100%;
}

/* 撒花动画 */
.confetti-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  pointer-events: none;
  z-index: 9999;
  overflow: hidden;
}

.confetti {
  position: absolute;
  top: -10px;
  width: 10px;
  height: 10px;
  border-radius: 2px;
  animation: confetti-fall linear forwards;
}

@keyframes confetti-fall {
  0% {
    transform: translateY(0) rotate(0deg) scale(1);
    opacity: 1;
  }
  100% {
    transform: translateY(100vh) rotate(720deg) scale(0.5);
    opacity: 0;
  }
}
</style>
