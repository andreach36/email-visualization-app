<script setup lang="ts">
import SearchBar from '@/components/SearchBar.vue';
import { useEmailStore } from '@/stores/data';
import type { Email } from '@/types';
import { ref, onMounted } from 'vue';

const emailStore = useEmailStore()
const isModalOpen = ref(false);

// Mostrar contenido del email
function showBody(emailData: Email){
  emailStore.emailVisualization.Subject = emailData.Subject
  emailStore.emailVisualization.Body = emailData.Body
  isModalOpen.value = true;
}

// Cerrar el modal
function closeModal() {
  isModalOpen.value = false;
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
      <!-- Detalles del email (Desktop) -->
      <div class="hidden md:block bg-white px-6 py-4 m-4 md:m-0 md:mr-8 rounded-lg h-[70vh] overflow-auto">
        <h2 class="font-bold text-lg">{{ emailStore.emailVisualization.Subject }}</h2>
        <p class="tracking-normal leading-relaxed text-justify text-sm py-3 overflow-hidden">{{ emailStore.emailVisualization.Body }}</p>
      </div>
    </div>

    <!-- Modal (Mobile) -->
    <div v-if="isModalOpen" class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50 md:hidden">
      <div class="bg-white rounded-lg p-4 max-w-sm w-full max-h-[80vh] overflow-auto relative">
        <button @click="closeModal" class="absolute top-2 right-2 bg-red-700 rounded text-white text-2xl font-bold px-1">&times;</button>
        <h2 class="font-bold text-lg">{{ emailStore.emailVisualization.Subject }}</h2>
        <p class="tracking-normal leading-relaxed text-justify text-sm py-3">{{ emailStore.emailVisualization.Body }}</p>
      </div>
    </div>

  </main>
</template>

