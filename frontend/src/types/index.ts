export interface Email {
    Id: number
    From: string
    To: string
    Subject: string
    Body: string
}

export interface DataPagination {
    page_num: number
    page_size: number
    total_pages: number
    total_data: number
}

export interface EmailVisualization {
    Subject: string
    Body: string
}