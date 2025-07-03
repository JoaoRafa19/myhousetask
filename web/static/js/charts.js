// Aguarda o documento inteiro carregar antes de rodar o script
document.addEventListener('DOMContentLoaded', function() {
    
    // Função para buscar os dados do nosso novo endpoint
    async function fetchChartData() {
        try {
            const response = await fetch('/api/charts/weekly-activity');
            if (!response.ok) {
                throw new Error('Falha ao buscar dados do gráfico');
            }
            return await response.json();
        } catch (error) {
            console.error('Erro no gráfico:', error);
            // Retorna dados padrão em caso de erro para não quebrar a página
            return { labels: [], data: [] };
        }
    }

    // Função para renderizar o gráfico com os dados
    function renderChart(chartData) {
        const ctx = document.getElementById('activityChartCanvas');
        if (!ctx) return; // Não faz nada se o canvas não existir na página

        new Chart(ctx, {
            type: 'bar',
            data: {
                labels: chartData.labels, // Ex: ['Seg', 'Ter', 'Qua', ...]
                datasets: [{
                    label: 'Tarefas Concluídas',
                    data: chartData.data, // Ex: [5, 9, 3, ...]
                    backgroundColor: 'rgba(79, 70, 229, 0.8)', // Cor índigo do Tailwind
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
                            // Garante que o eixo Y só mostre números inteiros
                            stepSize: 1
                        }
                    }
                },
                plugins: {
                    legend: {
                        display: false // Esconde a legenda para um look mais limpo
                    }
                }
            }
        });
    }

    // Executa o processo: busca os dados e depois renderiza o gráfico
    fetchChartData().then(renderChart);
});