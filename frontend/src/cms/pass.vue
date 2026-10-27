<template>
  <div class="cms-root">
    <router-view />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import AOS from 'aos'
import 'aos/dist/aos.css'
import '@/cms/styles/main.css'

onMounted(() => {
  AOS.init({ duration: 800, once: true, offset: 100 })
})
</script>
