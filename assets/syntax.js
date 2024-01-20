Prism.languages.goscript = {
    'comment': /\/\/.*/,
    'string': /"(?:\\.|[^\\"])*"/,
    'function': {
        pattern: /(\bfunc\s+)[a-zA-Z_]\w*(?=\()/,
        lookbehind: true
    },
    'keyword': /\b(?:if|for|return)\b/,
    'boolean': /\b(?:true|false)\b/,
    'number': /\b\d+(?:\.\d+)?\b/,
    'operator': /=/,
    'punctuation': /[{}[\];(),.:]/
};