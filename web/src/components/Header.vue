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

        <div id="dropdown-examples" class="dropdown is-active">
          <div class="dropdown-trigger">
            <button class="button" aria-haspopup="true" aria-controls="dropdown-menu2">
              <span>Examples</span>
              <img class="button-icon" src="../assets/arrow-down-icon.png">
            </button>
          </div>

          <div class="dropdown-menu" role="menu">
            <div class="dropdown-content panel" style="max-height: 229px;">
              <div class="panel-items"><div class="dropdown-item dropdown-item-section has-text-weight-semibold">
          Access Control
        </div><a class="dropdown-item" id="access-control.1-rbac">
          <div>Role-based</div>
          <div class="is-size-7 has-text-grey policy-description" style="width: 468px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
            An example of classic Role-based Access Control (RBAC)
          </div>
        </a><a class="dropdown-item" id="access-control.2-abac">
          <div>Attribute-based</div>
          <div class="is-size-7 has-text-grey policy-description" style="width: 468px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
            An example of classic Attribute-based Access Control (ABAC)
          </div>
        </a><div class="dropdown-item dropdown-item-section has-text-weight-semibold">
          Envoy
        </div><a class="dropdown-item" id="envoy.1-hello-world">
          <div>Hello World</div>
          <div class="is-size-7 has-text-grey policy-description" style="width: 468px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
            Allow everyone access to public APIs
          </div>
        </a><a class="dropdown-item" id="envoy.2-jwt">
          <div>JWT Decoding</div>
          <div class="is-size-7 has-text-grey policy-description" style="width: 468px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
            Verify a JWT and use the claims for authorization
          </div>
        </a><a class="dropdown-item" id="envoy.3-roles">
          <div>Roles</div>
          <div class="is-size-7 has-text-grey policy-description" style="width: 468px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
            Allow access to an API based on whether the user is an admin
          </div>
        </a><a class="dropdown-item" id="envoy.4-urlextract">
          <div>URL Extraction</div>
          <div class="is-size-7 has-text-grey policy-description" style="width: 468px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
            Grant the user access to her own user profiles
          </div>
        </a><div class="dropdown-item dropdown-item-section has-text-weight-semibold">
          Kubernetes
        </div><a class="dropdown-item" id="kubernetes.1-hello-world">
          <div>Hello World</div>
          <div class="is-size-7 has-text-grey policy-description" style="width: 468px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
            Ensure every resource has a 'costcenter' label in the appropriate format
          </div>
        </a><a class="dropdown-item has-tooltip-top has-tooltip-multiline" id="kubernetes.2-existence" data-tooltip="
            Ensure every resource has a costcenter label and that the costcenter value is in the correct format
          ">
          <div>Label Existence</div>
          <div class="is-size-7 has-text-grey policy-description" style="width: 468px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
            Ensure every resource has a costcenter label and that the costcenter value is in the correct format
          </div>
        </a><a class="dropdown-item" id="kubernetes.3-images">
          <div>Image Safety</div>
          <div class="is-size-7 has-text-grey policy-description" style="width: 468px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
            Ensure every image in a pod comes from a trusted registry
          </div>
        </a><a class="dropdown-item" id="kubernetes.4-ingress">
          <div>Ingress Conflicts</div>
          <div class="is-size-7 has-text-grey policy-description" style="width: 468px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
            Ensure no ingress gets created that conflicts with an existing ingress
          </div>
        </a></div>
            </div>
          </div>
        </div>
      </div>

      <div class="buttons level-right toolbar-item">
        <button @click="evaluateOnClick" id="evaluate-button" class="level-item button action-button is-info has-tooltip-bottom"
          data-tooltip="Query & view results in the output">
          <img class="button-icon" src="../assets/play-icon.png" />
          <span id="evaluate-button-text" class="button-label">Evaluate</span>
        </button>

        <button @click="formatOnClick" id="format-button" class="level-item button action-button is-info has-tooltip-bottom"
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
    margin-right: 0.75rem;
}

#logo {
    height: 100%;
    padding: 5px;
    margin-left: 0px;
    margin-right: 0px;
    box-sizing: border-box;
}

.dropdown.is-active .dropdown-menu, .dropdown.is-hoverable:hover .dropdown-menu {
    display: block;
}
.dropdown {
    display: inline-flex;
    position: relative;
    vertical-align: top;
}
.dropdown-trigger > button {
    justify-content: flex-start;
}

.dropdown-trigger > button {
    background-color: #fff;
    border-color: #dbdbdb;
    border-width: 1px;
    color: #363636;
    cursor: pointer;
    justify-content: center;
    padding-bottom: calc(0.5em - 1px);
    padding-left: 1em;
    padding-right: 1em;
    padding-top: calc(0.5em - 1px);
    text-align: center;
    white-space: nowrap;
}

.dropdown-trigger > button:hover {
    border-color: #b5b5b5;
    color: #363636;
}

.dropdown-menu {
    display: none;
    left: 0;
    min-width: 12rem;
    padding-top: 4px;
    position: absolute;
    top: 100%;
    z-index: 20;
}

.dropdown-content {
    min-width: 500px;
    overflow-y: auto;
    margin-top: 6px;
}

.dropdown > .button, .file-cta, .file-name, .input, .pagination-ellipsis, .pagination-link, .pagination-next, .pagination-previous, .select select, .textarea {
    -moz-appearance: none;
    -webkit-appearance: none;
    align-items: center;
    border: 1px solid transparent;
    border-radius: 4px;
    box-shadow: none;
    display: inline-flex;
    font-size: 1rem;
    height: 2.5em;
    justify-content: flex-start;
    line-height: 1.5;
    padding-bottom: calc(0.5em - 1px);
    padding-left: calc(0.75em - 1px);
    padding-right: calc(0.75em - 1px);
    padding-top: calc(0.5em - 1px);
    position: relative;
    vertical-align: top;
}

.dropdown-content {
    background-color: #242424;
    border-radius: 4px;
    box-shadow: 0 0.5em 1em -0.125em rgba(10,10,10,.1), 0 0 0 1px rgba(10,10,10,.02);
    padding-bottom: 0.5rem;
    padding-top: 0.5rem;
}

@media (prefers-color-scheme: light) {
    .dropdown-content {
        background-color: #fff;
    }
}

div.dropdown-item-section:first-child {
    border-top: 0;
}

div.dropdown-item-section {
    border-top: 1px solid #dbdbdb;
    padding: 6px 10px;
    cursor: default;
}

.dropdown-item {
    color: #4a4a4a;
    display: block;
    font-size: .875rem;
    line-height: 1.5;
    padding: 0.375rem 1rem;
    position: relative;
}

a.dropdown-item:hover, button.dropdown-item:hover {
    background-color: #f5f5f5;
    color: #0a0a0a;
}

.has-text-weight-semibold {
    font-weight: 600!important;
}

.has-text-grey {
    color: #7a7a7a!important;
}

.is-size-7 {
    font-size: .75rem!important;
}
</style>