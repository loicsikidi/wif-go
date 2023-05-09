<script setup>
    import { state } from '../js/state'
    import { utils } from '../js/utils'

    const props = defineProps({
        title: {
            type: String,
            required: true
        }
    })

    const evaluateOnClick = async () => {
        let output = ''
        let isErr = false
        try {
            if (!utils.isValidJson(state.get('mapping'))) {
                throw "ERROR: Input given to the compiler is not valid JSON..."
            }
            const obj = JSON.parse(state.get('mapping'))
            output = await wif_run(state.get('input'), obj.mapping || {}, obj.condition || '')
        } catch (error) {
            isErr = true
            console.error(error)
            if (typeof error === 'string') {
                output = error
            } else {
                output = error.message
            }
        } finally {
            state.update('output', output, !isErr)
        }
    }

    const formatOnClick = async () => {
        for (const key in state.getCode()){
            if (!utils.isValidJson(state.get(key))) {
                continue
            }
            state.update(key, JSON.parse(state.get(key)))
        }
    }
</script>

<template>
    <div id="header" class="level toolbar">
      <div class="level-left toolbar-item">
        <img src="/icon.png" id="logo" alt="Logo" class="level-item" />

        <h3 id="title" class="level-item title is-4">{{ title }}</h3>
      </div>

      <div class="buttons level-right toolbar-item">
        <button @click="evaluateOnClick" id="evaluate-button" class="level-item button is-info has-tooltip-bottom"
          data-tooltip="Query & view results in the output">
          <img class="button-icon" src="../assets/play-icon.png" />
          <span id="evaluate-button-text" class="button-label">Evaluate</span>
        </button>

        <button @click="formatOnClick" id="format-button" class="level-item button is-info has-tooltip-bottom"
          data-tooltip="Reformat the policy.">
          <img class="button-icon" src="../assets/format-icon.png" />
          <span class="button-label">Format</span>
        </button>
      </div>
    </div>
</template>

<style scoped>
#header {
    margin: 0;
    height: 56px;
}

#title {
    margin-bottom: 0;
    margin-top: 0;
}

#logo {
    height: 100%;
    padding: 5px;
    margin-left: 0px;
    margin-right: 0px;
    box-sizing: border-box;
}

</style>