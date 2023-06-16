<template>
  <div id="container">
      <div id="page" @mousemove="onDrag" @mouseup="endDrag">   
        <div id="policy-editor">
          <Editor
            :name="'mapping'"
            :config="configs['mapping'].editor_config"
            :language="configs['mapping'].language"
            :lint="configs['mapping'].lint"
            :defaultValue="configs['mapping'].defaultValue"
          />
        </div>
        <div class="gutter gutter-horizontal" @mousedown="startHzDrag"></div>
        <div id="input">
          <Editor
            :name="'input'"
            :config="configs['input'].editor_config"
            :language="configs['input'].language"
            :lint="configs['mapping'].lint"
            :defaultValue="configs['input'].defaultValue"
          />
        </div>
        <div class="gutter gutter-vertical" @mousedown="startVtDrag"></div>
        <div id="output">
          <Editor
            :name="'output'"
            :config="configs['output'].editor_config"
            :language="configs['output'].language"
            :defaultValue="configs['output'].defaultValue"
          />
        </div>
      </div>
  </div>
</template>
  
<script setup>
    import { computed, ref, toRaw, onMounted } from 'vue'
    import { json } from '@codemirror/lang-json'
    import { jsonParseLinter } from '@codemirror/lang-json'
    import Editor from './Editor.vue'

    const props = defineProps({
      editors: {
        type: Object,
        required: true
      }
    })

  const configs = computed(() => {
    let configs = {}
    const _defaultConfig = {
        disabled: false,
        indentWithTab: true,
        tabSize: 2,
        autofocus: true,
        height: 'calc(100% - 1.5rem)'
    }

    for (const [key, config] of Object.entries(toRaw(props.editors))) {
      const {header, placeholder} = config
      configs[key] = {
        editor_config: {...{header, placeholder}, ..._defaultConfig},
        defaultValue: config.defaultValue,
        language: json,
        lint: jsonParseLinter
      }
    }
    return configs
  })

    const hzDragging = ref(false)
    const vtDragging = ref(false)

    const startHzDrag = () => {
      hzDragging.value = true
    }

    const startVtDrag = () => {
      vtDragging.value = true
    }

    const endDrag = () => {
	    hzDragging.value = false
      vtDragging.value = false
    }

    const onDrag = (event) => {
      if(hzDragging.value || vtDragging.value) {
        const page = document.getElementById("page")
        const dragbarSize = 4
        const policyEditorWidth = event.clientX
        const upRowHeight = event.clientY

        if (hzDragging.value){
          const rightColWidth = page.clientWidth - dragbarSize - policyEditorWidth        
          const newColDefn = [policyEditorWidth, dragbarSize, rightColWidth]
            .map(c => c.toString() + "px").join(" ")
          page.style.gridTemplateColumns = newColDefn
        } else {
          const downRowHeight = page.clientHeight - dragbarSize - upRowHeight        
          const newRowDefn = [upRowHeight, dragbarSize,downRowHeight]
            .map(r => r.toString() + "px").join(" ")
          page.style.gridTemplateRows = newRowDefn
        }     
        event.preventDefault()
      }
    }

    onMounted(() => {
      // [HACK]: emulate resize event to set fixed width of the editor
      const gutter = document.getElementsByClassName("gutter-horizontal")[0]
      const coordonates = gutter.getBoundingClientRect()
      startHzDrag()
      onDrag(new MouseEvent('mousemove', {
        clientX: coordonates.x,
        clientY: coordonates.y
      }))
      endDrag()
    })
  </script>

<style scoped>
#container {
  display: inherit;
}

#page {
  display: grid;
	grid-template-areas:
		'policy-editor gutter-horizontal input'
		'policy-editor gutter-horizontal gutter-vertical'
		'policy-editor gutter-horizontal output';
	grid-template-rows: 1fr 6px 1fr;
	grid-template-columns: 6fr 6px 4fr;
}

/*****************************/
#policy-editor {
	overflow: auto;
	grid-area: policy-editor;
}

#input {
	overflow: auto;
	grid-area: input;
}

#output {
	overflow: auto;
	grid-area: output;
}

.gutter {
  background-color: rgba(0, 0, 0, 0.115);
  background-repeat: no-repeat;
  background-position: 50%;
}
.gutter.gutter-horizontal {
  height: 100%;
  cursor: col-resize;
  width: 4px;
  grid-area: gutter-horizontal;
  background-image: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAeCAYAAADkftS9AAAAIklEQVQoU2M4c+bMfxAGAgYYmwGrIIiDjrELjpo5aiZeMwF+yNnOs5KSvgAAAABJRU5ErkJggg==);
}

.gutter.gutter-vertical {
  height: 4px;
  cursor: row-resize;
  grid-area: gutter-vertical;
  background-image: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAB4AAAAFAQMAAABo7865AAAABlBMVEVHcEzMzMzyAv2sAAAAAXRSTlMAQObYZgAAABBJREFUeF5jOAMEEAIEEFwAn3kMwcB6I2AAAAAASUVORK5CYII=);
}
</style>