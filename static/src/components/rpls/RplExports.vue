<template>
  <div class="uk-position-relative">
    <h3>Functions</h3>
    <div class="uk-overflow-auto">
      <table class="uk-table">
        <tr v-for="item in exports.functions">
          <td>
            <code>
              {{ item }}
            </code>
          </td>
        </tr>
      </table>
    </div>

    <h3>Data</h3>
    <div class="uk-overflow-auto">
      <table class="uk-table">
        <tr v-for="item in exports.data">
          <td>
            <code>
              {{ item }}
            </code>
          </td>
        </tr>
      </table>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';

export default {
  props: ['titleId', 'rplId'],

  computed: {
    ...mapGetters(['exports', 'loadingExports']),
  },

  methods: {
    ...mapActions(['listExports']),

    fetchExports() {
      const { titleId, rplId } = this;

      if (titleId && rplId) {
        this.listExports({ titleId, rplId });
      }
    },
  },

  watch: {
    rplId() {
      this.fetchExports();
    },
    titleId() {
      this.fetchExports();
    },
  },

  beforeMount() {
    this.fetchExports();
  },
};
</script>

<style>
</style>
