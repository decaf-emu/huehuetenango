<template>
  <div class="uk-container uk-container-expand">
    <ul class="uk-breadcrumb">
      <li><router-link :to="{ name: 'titles'}">Titles</router-link></li>
      <li><span>Import</span></li>
    </ul>

    <h1>Import</h1>

    <form @submit="formSubmit">
      <div class="uk-margin" uk-margin>
        <div class="uk-alert-success" uk-alert v-show="importSuccess">Import succeeded.</div>
        <div class="uk-alert-danger" uk-alert v-show="importFailed">Something went wrong.</div>
        <div uk-form-custom="target: true">
          <input type="file" @change="fileChange($event.target.files[0]);" :disabled="importing">
          <input class="uk-input uk-form-width-medium" type="text" placeholder="Select file" disabled>
        </div>
        <button class="uk-button uk-button-default" :disabled="file === null || importing">Submit</button>
        <div class="uk-margin-left" uk-margin uk-spinner v-show="importing"></div>
      </div>
    </form>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  computed: mapGetters({
    importing: 'importingTitles',
    importSuccess: 'importTitlesSuccess',
    importFailed: 'importTitlesFailed',
  }),
  data() {
    return {
      file: null,
    };
  },
  methods: {
    fileChange(file) {
      this.file = file;
    },
    formSubmit(e) {
      this.$store.dispatch('importTitles', this.file);
      e.preventDefault();
    },
  },
};
</script>

<style>
.import-container {
}
</style>
