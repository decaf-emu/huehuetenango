<template>
  <div class="uk-container uk-container-expand">
    <ul class="uk-breadcrumb">
      <li><router-link :to="{ name: 'titles'}">Titles</router-link></li>
      <li><span>Import</span></li>
    </ul>

    <h1>Import</h1>

    <form @submit="formSubmit">
      <div class="uk-margin" uk-margin>
        <div class="uk-alert-success" uk-alert v-show="importTitlesSuccess">Import succeeded.</div>
        <div class="uk-alert-danger" uk-alert v-show="importTitlesFailed">Something went wrong.</div>
        <div uk-form-custom="target: true">
          <input type="file" @change="fileChange($event.target.files[0])" :disabled="importingTitles">
          <input class="uk-input uk-form-width-medium" type="text" placeholder="Select file" disabled>
        </div>
        <button class="uk-button uk-button-default" :disabled="file === null || importingTitles">Submit</button>
        <div class="uk-margin-left" uk-margin uk-spinner v-show="importingTitles"></div>
      </div>
    </form>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';

export default {
  data() {
    return {
      file: null,
    };
  },

  computed: {
    ...mapGetters([
      'importingTitles',
      'importTitlesSuccess',
      'importTitlesFailed',
    ]),
  },

  methods: {
    ...mapActions(['importTitles']),

    fileChange(file) {
      this.file = file;
    },

    formSubmit(e) {
      e.preventDefault();
      this.importTitles(this.file);
    },
  },
};
</script>

<style>
.import-container {
}
</style>
