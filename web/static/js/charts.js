async function fetchChartData() {
    try {
        const response = await fetch('/api/charts/weekly-activity');
        if (!response.ok) {
            throw new Error('Falha ao buscar dados do gráfico');
        }
        return await response.json();
    } catch (error) {
        console.error('Erro no gráfico:', error);
        return { labels: [], data: [] };
    }
}

let myChart = null; // Variável global para guardar a instância do gráfico

function renderChart(chartData) {
    const ctx = document.getElementById('activityChartCanvas');
    if (!ctx) return;

    if (myChart) {
        myChart.destroy();
    }

    myChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: chartData.labels,
            datasets: [{
                label: 'Tarefas Concluídas',
                data: chartData.data,
                backgroundColor: 'rgba(79, 70, 229, 0.8)',
                borderColor: 'rgba(79, 70, 229, 1)',
                borderWidth: 1,
                borderRadius: 4,
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        stepSize: 1
                    }
                }
            },
            plugins: {
                legend: {
                    display: false
                }
            }
        }
    });
}

function initializeDashboardChart() {
    if (document.getElementById('activityChartCanvas')) {
        fetchChartData().then(renderChart);
    }
}

document.addEventListener('DOMContentLoaded', initializeDashboardChart);

// 2. Executa toda vez que o HTMX troca conteúdo na página
document.body.addEventListener('htmx:afterSwap', initializeDashboardChart);