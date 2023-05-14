<script setup>
import { reactive, shallowRef } from 'vue'
import Header from './components/Header.vue'
import Grid from './components/Grid.vue'
import Footer from './components/Footer.vue'
import Loading from './components/common/Loading.vue'

const editors = {
  input: {
    header: "input",
    placeholder: "Please enter the input token (eg. JWT).",
    defaultValue: JSON.parse(`{"aud":"awesome_demo","sub": "id/123456789", "workload_id":"55d36609-9bcf-48e0-a366-a3cf19027d2a", "groups":["admins", "editors"], "iat":${Math.floor(Date.now() / 1000)}, "exp":${Math.floor(Date.now() / 1000) + 3600}}`),
    language: 'json',
  },
  output: {
    header: "output",
    placeholder: null,
    language: 'json',
  },
  mapping: {
    header: "attribute mapping",
    placeholder: "Please enter a JSON payload.",
    defaultValue: JSON.parse(`{"mapping":{"google.subject": "'myprovider::' + assertion.aud + '::' + assertion.sub.extract('id/{id}')", "google.groups":"has(assertion.groups) ? assertion.groups : []", "attribute.workload_display_name": "{'8bb39bdb-1cc5-4447-b7db-a19e920eb111': 'Workload1', '55d36609-9bcf-48e0-a366-a3cf19027d2a': 'Workload2'}[assertion.workload_id]"},"condition":"assertion.sub.extract('id/{id}') == '123456789' && 'admins' in google.groups"}`),
    language: 'json',
  },
}

const version = shallowRef('unknown')
const loading = shallowRef(true)
const config = reactive({
  editors: editors,
})

window.addEventListener('DOMContentLoaded', async () => {
        try {
          console.log("Loading WIF-GO WASM module ⏳⏳⏳")
          await loadWif()
          console.log("Successfully loaded WIF-GO WASM module 🚀")
          version.value = await wif_version()
          loading.value = false
        } catch (e) {
          console.error('Failed to load WIF-GO WASM: ' + e)
          alert('Failed to load WIF-GO WASM: ' + e)
        }

        async function loadWif() {
          if (!WebAssembly.instantiateStreaming) { // polyfill
              WebAssembly.instantiateStreaming = async (resp, importObject) => {
                  const source = await (await resp).arrayBuffer()
                  return await WebAssembly.instantiate(source, importObject)
              }
          }

          function loadWasm(path) {
              const go = new Go()
              return new Promise((resolve, reject) => {
                  WebAssembly.instantiateStreaming(fetch(path), go.importObject)
                      .then(result => {
                          go.run(result.instance)
                          resolve(result.instance)
                      })
                      .catch(error => {
                          reject(error)
                      })
                  })
          }
          await loadWasm("/wif-go.wasm")
      }
})
</script>

<template>
  <Header title="Workload Identity Federation Playground" />
  <div class="loading-box" v-if="loading">
      <Loading />
  </div>
  <div id="body" v-else>
    <div id="playground-content">
      <Grid :editors="config.editors"  />
    </div>
    <Footer :version="version" />
  </div>
</template>

<style lang="css">
  #body {
      position: absolute;
      top: 56px;
      right: 0;
      bottom: 0;
      left: 0;
      display: flex;
      flex-direction: column;
  }

  #playground-content {
    border: #0a0a0a;
    overflow: hidden;
    display: flex;
    flex-grow: 1;
    position: relative;
  }

  .loading-box {
    width: 100%;
    min-height: 20rem;
    max-height: 60rem;
    padding-top: 20%;
  }
</style>
