import { reactive } from 'vue'
import { utils } from './utils'

export const state = reactive({
    code: {
        mapping: "",
        input: "",
        output: "",
    },
    getCode() {
      return this.code
    },
    get(key) {
      return this.code[key]
    },
    update(key, value, isJson = true) {
      this.code[key] = isJson ? utils.prettyJson(value) : value
    }
})