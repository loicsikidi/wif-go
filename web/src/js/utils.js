export const utils = {
    prettyJson(str) {
        try {
            return JSON.stringify(str, null, 2)
        } catch (e) {
            return str
        }
    },
    isValidJson(str) {
        try {
            JSON.parse(str)
            return true
        } catch(e) {
            return false
        }
    },
    sleep(ms){
        return new Promise(r => setTimeout(r, ms))
    }
}