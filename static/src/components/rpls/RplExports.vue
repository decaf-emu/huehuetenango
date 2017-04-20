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
import { mapGetters } from 'vuex';

export default {
  computed: {
    ...mapGetters(['exports']),
  },
};
</script>

<style>
</style>
