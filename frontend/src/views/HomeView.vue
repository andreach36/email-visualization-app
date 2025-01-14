<script setup lang="ts">
import SearchBar from '@/components/SearchBar.vue';
import { useEmailStore } from '@/stores/data';
import type { Email } from '@/types';
import { onMounted } from 'vue';

const emailStore = useEmailStore()

// Mostrar contenido del email
function showBody(emailData: Email){
  emailStore.emailVisualization.Subject = emailData.Subject
  emailStore.emailVisualization.Body = emailData.Body
}

// Cambiar página
function changePage(newPage: number) {
  if (newPage >= 1 && newPage <= emailStore.dataPagination.total_pages) {
    emailStore.loadData(newPage).then(() => {
      emailStore.dataPagination.page_num = newPage
    })
  }
}

onMounted(async () => {
  emailStore.loadData(emailStore.dataPagination.page_num) // cargar la data
})

</script>

<template>
  <main class="overflow-hidden bg-gray-100 py-16 min-h-screen">
    <SearchBar></SearchBar>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- Tabla -->
      <div class="flex flex-col px-6">
        <table class="table-fixed w-full border-collapse border border-gray-300">
          <thead>
            <tr class="bg-violet-300">
              <th class="border border-gray-300 px-1 md:px-4 py-2 text-xs font-medium text-gray-700">Subject</th>
              <th class="border border-gray-300 px-1 md:px-4 py-2 text-xs font-medium text-gray-700">From</th>
              <th class="border border-gray-300 px-1 md:px-4 py-2 text-xs font-medium text-gray-700">To</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="datarow in emailStore.dataEmails" :key="datarow.Id" class="bg-white">
              <td class="cursor-pointer border border-gray-300 px-1 md:px-4 py-2 text-center text-xs truncate text-ellipsis whitespace-nowrap" @click="showBody(datarow)">
                {{ datarow.Subject }}
              </td>
              <td class="cursor-pointer border border-gray-300 px-1 md:px-4 py-2 text-center text-xs truncate text-ellipsis whitespace-nowrap" @click="showBody(datarow)">
                {{ datarow.From }}
              </td>
              <td class="cursor-pointer border border-gray-300 px-1 md:px-4 py-2 text-center text-xs truncate text-ellipsis whitespace-nowrap" @click="showBody(datarow)">
                {{ datarow.To }}
              </td>
            </tr>
          </tbody>
        </table>
        <div class="flex justify-between items-center px-4 py-2 bg-gray-100 border-t">
          <button class="px-2 py-1 bg-sky-800 text-white text-sm rounded disabled:opacity-50"
          :disabled="emailStore.dataPagination.page_num <= 1"
          @click="changePage(emailStore.dataPagination.page_num - 1)">Before</button>
          <span class="text-sm text-center text-gray-700">Page {{ emailStore.dataPagination.page_num }} from {{ emailStore.dataPagination.total_pages }}</span>
          <button class="px-2 py-1 bg-sky-800 text-white text-sm rounded disabled:opacity-50"
          :disabled="emailStore.dataPagination.page_num >= emailStore.dataPagination.total_pages"
          @click="changePage(emailStore.dataPagination.page_num + 1)">Next</button>
          <span class="text-sm text-center text-gray-700">Total Data: {{ emailStore.dataPagination.total_data }} emails</span>
        </div>
      </div>
      <!-- Detalles del email -->
      <div class="bg-white px-6 py-4 m-4 md:m-0 md:mr-8 rounded-lg h-[70vh] overflow-auto">
        <h2 class="font-bold text-lg">{{ emailStore.emailVisualization.Subject }}</h2>
        <p class="tracking-normal leading-relaxed text-justify text-sm py-3 overflow-hidden">{{ emailStore.emailVisualization.Body }}</p>
      </div>
    </div>
  </main>
</template>

