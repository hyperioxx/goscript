<link href="https://cdn.jsdelivr.net/npm/prismjs@1.25.0/themes/prism.css" rel="stylesheet"/>
<script src="https://cdn.jsdelivr.net/npm/prismjs@1.25.0/prism.js"></script>
<script>
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
</script>

# GoScript

GoScript is a dynamically typed, interpreted language created out of curiosity to answer that question we ask as programmers: "How do you make a programming language from scratch?" So, I've given it a try.


Example syntax:
<pre><code class="language-goscript">
// variable declaration
myInt = 1
myString = "foo"
myFloat = 1.0
myArray = [1,2,3,4,"bar"]



// conditionals
if myInt > 1 {
    print(x) // builtin function 
}
</code></pre>

