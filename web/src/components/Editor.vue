<template>
    <div class="panel-title" v-if="config.header">
        <span>{{ config.header.toUpperCase() }}</span>
    </div>
    <codemirror
      v-model="state.code[props.name]"
      :style="{
        width: '100%',
        height: config.height,
        backgroundColor: '#fff',
        color: '#333'
      }"
      :placeholder="config.placeholder"
      :extensions="extensions"
      :autofocus="config.autofocus"
      :disabled="config.disabled"
      :indent-with-tab="config.indentWithTab"
      :tab-size="config.tabSize"
      @ready="handleReady"
    />
</template>
  
<script setup>
    import { shallowRef, computed } from 'vue'
    import { EditorView } from '@codemirror/view'
    import { linter } from '@codemirror/lint'
    import { Codemirror } from 'vue-codemirror'
    import { state } from '../js/state'

    const props = defineProps({
      name: {
        type: String,
        required: true
      },
      config: {
        type: Object,
        required: true
      },
      defaultValue: {
        type:[String, Object],
        required: false
      },
      theme: [Object, Array],
      language: Function,
      lint: Function
    })

    // [INFO]: set default value to state
    if (props.defaultValue)
      state.update(props.name, props.defaultValue)

    const extensions = computed(() => {
        const result = []
        if (props.language) {
          result.push(props.language())
        }
        if (props.lint){
          result.push(linter(props.lint()))
        }
        if (props.theme) {
          result.push(props.theme)
        }
        return result
    })    
    const cmView = shallowRef(EditorView)
    const handleReady = (payload) => {
        cmView.value = payload.view
    }
</script>

<style scoped>
.panel-title {
    display: flex;
    background: #999999;
    color: white;
    height: 1.5rem;
    font-size: 1rem;
    font-weight: 500;
    margin: 0;
    padding-left: 5px;
    user-select: none;
}
</style>