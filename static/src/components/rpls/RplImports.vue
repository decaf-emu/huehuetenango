<template>
  <div class="uk-position-relative">
    <h3>Functions</h3>
    <div class="uk-overflow-auto">
      <table class="uk-table">
        <RplImportRow v-for="item in imports.functions" :item="item" :key="item.name" />
      </table>
    </div>

    <h3>Data</h3>
    <div class="uk-overflow-auto">
      <table class="uk-table">
        <RplImportRow v-for="item in imports.data" :item="item" :key="item.name" />
      </table>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import RplImportRow from './RplImportRow.vue';

export default {
  components: { RplImportRow },
  props: ['titleId', 'rplId'],

  computed: {
    ...mapGetters(['imports', 'loadingImports']),
  },

  methods: {
    ...mapActions(['listImports']),

    fetchImports() {
      const { titleId, rplId } = this;

      if (titleId && rplId) {
        this.listImports({ titleId, rplId });
      }
    },
  },

  watch: {
    rplId() {
      this.fetchImports();
    },

    titleId() {
      this.fetchImports();
    },
  },

  beforeMount() {
    this.fetchImports();
  },
};
</script>

<style>
</style>
