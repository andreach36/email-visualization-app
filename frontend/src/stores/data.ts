import { ref} from 'vue'
import { defineStore } from 'pinia'
import type { Email } from '@/types'
import type { DataPagination } from '../types/index';
import type { EmailVisualization } from '../types/index';

export const useEmailStore = defineStore('email', () => {
  
  const dataEmails = ref<Email[]>([])

  const dataPagination = ref<DataPagination>({
    page_num: 1,
    page_size: 10,
    total_pages: 0,
    total_data: 0,
  })

  const searchQuery = ref<string>('')

  const emailVisualization = ref<EmailVisualization>({
    Subject: "",
    Body: "",
  })


  async function loadData (page: number) {
    try {
      let baseUrl = import.meta.env.VITE_API_URL
      let paramsString = '';

      if (searchQuery.value === '') {
        paramsString = new URLSearchParams({
          page_num: String(page)
        }).toString();
      } else {
        baseUrl = baseUrl + "/search"
        paramsString = new URLSearchParams({
          q: searchQuery.value,
          page_num: String(page)
        }).toString();
      }

      const url = `${baseUrl}?${paramsString}`;

      // Realizar la solicitud
      const response = await fetch(url, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        },
      })
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`)
      }
      const result = await response.json()
      
      // Mapear los datos emails
      dataEmails.value = result.data.map((email: any) => ({
        Id: email.ID,
        From: email.from,
        To: email.to,
        Subject: email.subject,
        Body: email.Body,
      }));

      // Obtener datos paginación
      dataPagination.value = (result.pagination as DataPagination)

      // Obtener data panel visualizador
      emailVisualization.value.Subject = dataEmails.value[0].Subject
      emailVisualization.value.Body = dataEmails.value[0].Body
      
    } catch (error) {
      console.error('Failed to load data:', error)
      dataEmails.value = []
    }
  }

  return {dataEmails, dataPagination, searchQuery, emailVisualization, loadData}

})