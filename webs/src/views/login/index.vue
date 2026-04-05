<template>
  <div class="login-container">
    <div class="login-bg-shapes">
      <div class="shape shape-1"></div>
      <div class="shape shape-2"></div>
      <div class="shape shape-3"></div>
    </div>

    <div class="absolute-lt flex-x-end p-3 w-full" style="z-index: 10;">
      <el-switch
        v-model="isDark"
        inline-prompt
        :active-icon="Moon"
        :inactive-icon="Sunny"
        @change="toggleTheme"
      />
      <lang-select class="ml-2 cursor-pointer" />
    </div>

    <el-card class="login-card !border-none !rounded-4% w-100 <sm:w-85">
      <div class="login-header">
        <div class="login-logo">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" class="logo-icon">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" fill="currentColor" opacity="0.8"/>
            <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <h2 class="login-title">{{ defaultSettings.title }}</h2>
        <el-tag class="version-tag" effect="plain" round>{{ version }}</el-tag>
      </div>

      <el-form
        ref="loginFormRef"
        :model="loginData"
        :rules="loginRules"
        class="login-form"
      >
        <el-form-item prop="username">
          <el-input
            ref="username"
            v-model="loginData.username"
            :placeholder="$t('login.username')"
            name="username"
            size="large"
            :prefix-icon="User"
            class="login-input"
          />
        </el-form-item>

        <el-tooltip
          :visible="isCapslock"
          :content="$t('login.capsLock')"
          placement="right"
        >
          <el-form-item prop="password">
            <el-input
              v-model="loginData.password"
              :placeholder="$t('login.password')"
              type="password"
              name="password"
              @keyup="checkCapslock"
              @keyup.enter="handleLogin"
              size="large"
              :prefix-icon="Lock"
              show-password
              class="login-input"
            />
          </el-form-item>
        </el-tooltip>

        <el-form-item prop="captchaCode">
          <div class="captcha-row">
            <el-input
              v-model="loginData.captchaCode"
              auto-complete="off"
              size="large"
              :prefix-icon="Key"
              class="captcha-input"
              :placeholder="$t('login.captchaCode')"
              @keyup.enter="handleLogin"
            />
            <el-image
              @click="getCaptcha"
              :src="captchaBase64"
              class="captcha-image cursor-pointer"
            />
          </div>
        </el-form-item>

        <el-button
          :loading="loading"
          type="primary"
          size="large"
          class="login-btn"
          @click.prevent="handleLogin"
        >
          {{ $t("login.login") }}
        </el-button>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { useSettingsStore, useUserStore } from "@/store";
import { getCaptchaApi, GetVersion } from "@/api/auth";
import { LoginData } from "@/api/auth/types";
import { Sunny, Moon, User, Lock, Key } from "@element-plus/icons-vue";
import { LocationQuery, LocationQueryValue, useRoute } from "vue-router";
import router from "@/router";
import defaultSettings from "@/settings";
import { ThemeEnum } from "@/enums/ThemeEnum";

const version = ref('')
const fetchVersion = function(){
  GetVersion().then((res) => {
    version.value = res.data;
  }).catch((error) => {
    console.error("Error fetching version:", error);
  });
}()

const userStore = useUserStore();
const settingsStore = useSettingsStore();
const { t } = useI18n();

const isDark = ref(settingsStore.theme === ThemeEnum.DARK);
const loading = ref(false);
const isCapslock = ref(false);
const captchaBase64 = ref();
const loginFormRef = ref(ElForm);
const { height } = useWindowSize();

const loginData = ref<LoginData>({
  username: "",
  password: "",
});

const loginRules = computed(() => {
  return {
    username: [
      {
        required: true,
        trigger: "blur",
        message: t("login.message.username.required"),
      },
    ],
    password: [
      {
        required: true,
        trigger: "blur",
        message: t("login.message.password.required"),
      },
      {
        min: 6,
        message: t("login.message.password.min"),
        trigger: "blur",
      },
    ],
    captchaCode: [
      {
        required: true,
        trigger: "blur",
        message: t("login.message.captchaCode.required"),
      },
    ],
  };
});

function getCaptcha() {
  getCaptchaApi().then(({ data }) => {
    loginData.value.captchaKey = data.captchaKey;
    captchaBase64.value = data.captchaBase64;
  });
}

const route = useRoute();
function handleLogin() {
  loginFormRef.value.validate((valid: boolean) => {
    if (valid) {
      loading.value = true;
      userStore
        .login(loginData.value)
        .then(() => {
          const query: LocationQuery = route.query;
          const redirect = (query.redirect as LocationQueryValue) ?? "/";
          const otherQueryParams = Object.keys(query).reduce(
            (acc: any, cur: string) => {
              if (cur !== "redirect") {
                acc[cur] = query[cur];
              }
              return acc;
            },
            {}
          );
          router.push({ path: redirect, query: otherQueryParams });
        })
        .catch(() => {
          getCaptcha();
        })
        .finally(() => {
          loading.value = false;
        });
    }
  });
}

/**
 * 主题切换
 */

const toggleTheme = () => {
  const newTheme =
    settingsStore.theme === ThemeEnum.DARK ? ThemeEnum.LIGHT : ThemeEnum.DARK;
  settingsStore.changeTheme(newTheme);
};

function checkCapslock(event: KeyboardEvent) {
  if (event instanceof KeyboardEvent) {
    isCapslock.value = event.getModifierState("CapsLock");
  }
}

onMounted(() => {
  getCaptcha();
});
</script>

<style lang="scss" scoped>
.login-container {
  overflow-y: auto;
  position: relative;
  background: linear-gradient(135deg, #0F172A 0%, #1E293B 50%, #0F172A 100%);
  @apply wh-full flex-center;
}

html.dark .login-container {
  background: linear-gradient(135deg, #020617 0%, #0F172A 50%, #020617 100%);
}

.login-bg-shapes {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;

  .shape {
    position: absolute;
    border-radius: 50%;
    filter: blur(80px);
    opacity: 0.4;
  }

  .shape-1 {
    width: 400px;
    height: 400px;
    background: #2563EB;
    top: -100px;
    right: -100px;
  }

  .shape-2 {
    width: 300px;
    height: 300px;
    background: #6366F1;
    bottom: -80px;
    left: -80px;
  }

  .shape-3 {
    width: 200px;
    height: 200px;
    background: #F97316;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    opacity: 0.15;
  }
}

.login-card {
  position: relative;
  z-index: 5;
  background: rgba(255, 255, 255, 0.08) !important;
  backdrop-filter: blur(24px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.12) !important;
  box-shadow: 0 25px 50px -12px rgb(0 0 0 / 0.5) !important;
}

html:not(.dark) .login-card {
  background: rgba(255, 255, 255, 0.85) !important;
  border: 1px solid rgba(255, 255, 255, 0.5) !important;
  box-shadow: 0 25px 50px -12px rgb(0 0 0 / 0.15) !important;
}

.login-header {
  text-align: center;
  margin-bottom: 8px;
  position: relative;
}

.login-logo {
  display: flex;
  justify-content: center;
  margin-bottom: 12px;

  .logo-icon {
    width: 48px;
    height: 48px;
    color: #2563EB;
  }
}

.login-title {
  font-size: 24px;
  font-weight: 700;
  color: #F1F5F9;
  margin: 0 0 4px 0;
  letter-spacing: -0.5px;
}

html:not(.dark) .login-title {
  color: #1E293B;
}

.version-tag {
  position: absolute;
  top: 0;
  right: 0;
  font-size: 11px;
}

.login-form {
  padding: 16px 10px 0;
}

.login-input {
  :deep(.el-input__wrapper) {
    padding: 4px 12px;
    background-color: rgba(255, 255, 255, 0.06);
    border-radius: 10px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    box-shadow: none !important;
    transition: all 200ms ease-out;

    &:hover {
      border-color: rgba(37, 99, 235, 0.4);
    }

    &.is-focus {
      border-color: #2563EB;
      background-color: rgba(255, 255, 255, 0.1);
      box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.15) !important;
    }

    input {
      color: #F1F5F9;

      &:-webkit-autofill {
        transition: background-color 1000s ease-in-out 0s;
        -webkit-text-fill-color: #F1F5F9;
      }
    }
  }
}

html:not(.dark) .login-input {
  :deep(.el-input__wrapper) {
    background-color: rgba(241, 245, 249, 0.8);
    border-color: #E2E8F0;

    input {
      color: #1E293B;

      &:-webkit-autofill {
        -webkit-text-fill-color: #1E293B;
      }
    }
  }
}

.captcha-row {
  display: flex;
  width: 100%;
  gap: 8px;
}

.captcha-input {
  flex: 1;
}

.captcha-image {
  width: 120px;
  height: 40px;
  border-radius: 10px;
  overflow: hidden;
  flex-shrink: 0;
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 15px;
  font-weight: 600;
  border-radius: 10px !important;
  background: linear-gradient(135deg, #2563EB 0%, #3B82F6 100%) !important;
  border: none !important;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.35);
  transition: all 200ms ease-out !important;
  letter-spacing: 1px;

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 6px 20px rgba(37, 99, 235, 0.45) !important;
  }

  &:active {
    transform: translateY(0);
  }
}

.el-form-item {
  background: transparent;
  border: none;
  border-radius: 0;
  margin-bottom: 18px;
}
</style>
