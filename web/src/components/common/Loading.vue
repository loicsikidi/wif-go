<template>
    <div class="loading">
      <div class="animation">
        <div v-for="item in 5" :key="item"></div>
      </div>
    </div>
</template>
  
  <style lang="css" scoped>
    .loading {
      width: 100%;
      height: 100%;
      position: relative;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      user-select: none;
  
    }

    @keyframes animation {
          0%,
          100% {
            opacity: 1;
          }
          50% {
            opacity: 0.5;
          }
    }

      .animation {
        display: flex;
      }
  
      .animation div {
          width: 18px;
          height: 18px;
          margin: 0 9px;
          border-radius: 2px;
          background: #42b883;
          animation: animation 2s linear infinite;
        }

        .animation div:nth-child(1) {
          animation-delay: 0s;
        }
        .animation div:nth-child(2) {
          animation-delay: -0.4s;
        }
        .animation div:nth-child(3) {
          animation-delay: -0.8s;
        }
        .animation div:nth-child(4) {
          animation-delay: -0.4s;
        }
        .animation div:nth-child(5) {
          animation-delay: -0.8s;
        }
        .animation div:nth-child(6) {
          animation-delay: -1.2s;
        }
  </style>