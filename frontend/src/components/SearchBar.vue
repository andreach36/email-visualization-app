<script setup lang="ts">
import { useEmailStore } from '@/stores/data';
import { ref } from 'vue';

const emailStore = useEmailStore()
const searchInput = ref('');

function handleSearch(){
  emailStore.searchQuery = searchInput.value;
  emailStore.loadData(emailStore.dataPagination.page_num);

}

</script>

<template>
    <div class="flex gap-4 py-8 px-2">
      <input type="search" v-model="searchInput" placeholder="Search emails" class="px-2 py-2 border-2 border-slate-500 rounded text-sm w-4/5">
      <button class="bg-rose-400 px-2 py-2 rounded text-sm text-white font-medium w-1/5" @click="handleSearch">Search</button>
    </div>
</template>