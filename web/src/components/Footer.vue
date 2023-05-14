<script setup>
    const props = defineProps({
        version: {
            type: String,
            required: true
        }
    })
</script>

<template>
    <div id="footer" class="level toolbar">
        <div class="level-left toolbar-item">
          <!-- Note: Direct style attribute is a work-around for Safari.. -->
          <span class="is-size-7 has-text-weight-semibold" style="white-space: nowrap">Built with ❤️ by</span>
          <a id="github-logo-link" href="https://github.com/loicsikidi" target="_blank">
            <img src="../assets/github.svg" id="github-logo" alt="Github logo">
          </a>
        </div>
        <div class="level-right toolbar-item">wif-go&nbsp;<span class="wif-go-version"><a :href="`https://github.com/loicsikidi/google-auth-plugins-library-python/tree/${version}`" target="_blank">{{ version }}</a></span>
        </div>
    </div>
</template>

<style scoped>
#footer {
    height: 40px;
    font-size: .75rem;
    margin: 0;
    display: flex;
}

#github-logo-link {
    height: 60%;
}

#github-logo {
    height: 100%;
    object-fit: contain;
    padding-left: 5px;
    vertical-align: middle;
}
</style>