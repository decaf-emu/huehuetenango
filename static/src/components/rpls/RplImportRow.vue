<template>
  <tr>
    <td class="uk-table-expand uk-table-link" v-if="item.source_title_id && item.source_id">
      <router-link :to="{ name: 'title', params: { titleId: item.source_title_id, rplId: item.source_id, type: 'exports' }}">
        <code>{{ item.name }}</code>
      </router-link>
    </td>
    <td class="uk-table-expand" v-if="!item.source_title_id || !item.source_id">
      <code>{{ item.name }}</code>
    </td>

    <td class="uk-table-shrink uk-table-link" v-if="item.source_title_id && item.source_id">
      <router-link :to="{ name: 'title', params: { titleId: item.source_title_id, rplId: item.source_id, type: 'exports' }}">
        {{ item.source }}
      </router-link>
    </td>
    <td class="uk-table-shrink" v-if="!item.source_title_id || !item.source_id">
      {{ item.source }}
    </td>
  </tr>
</template>

<script>
export default {
  props: ['item'],
};
</script>

<style>
</style>
