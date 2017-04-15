# eslint-config-enough
The most useful eslint rules that just enough.

## Install

```
npm install eslint-config-enough --save-dev
```

## Usage:
Create `.eslintrc.js` at the root of your project:

```js
module.exports = {
  extends: 'enough',

  // es6 has been set by default
  env: {
    browser: true,
    node: true
  },

  // overrides
  "rules": {
    // ...
  }
}
```
