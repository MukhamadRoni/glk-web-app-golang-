document.addEventListener('DOMContentLoaded', function () {
    const API_CS = "/api/v1/confidence-score";
    const modalElement = document.getElementById('modalCS');
    const modalCS = new bootstrap.Modal(modalElement);
    const formCS = document.getElementById('formCS');
    const colorPicker = document.getElementById('csColorPicker');
    const colorInput = document.getElementById('csColor');

    let grid;

    // Sync color picker and hex input
    colorPicker.addEventListener('input', () => colorInput.value = colorPicker.value.toUpperCase());
    colorInput.addEventListener('change', () => colorPicker.value = colorInput.value);

    function initGrid() {
        grid = new gridjs.Grid({
            columns: [
                { name: "ID", hidden: true },
                { name: "Kategori", width: "200px", formatter: (cell, row) => {
                    const color = row.cells[2].data;
                    return gridjs.html(`<span class="badge" style="background-color: ${color}; color: ${getContrastYIQ(color)}; font-size: 0.9rem;">${cell}</span>`);
                }},
                { name: "Warna", width: "120px" },
                { name: "Rentang Skor", width: "150px", formatter: (cell, row) => `${row.cells[3].data} - ${row.cells[4].data}` },
                { name: "Min", hidden: true },
                { name: "Max", hidden: true },
                {
                    name: "Aksi",
                    width: "150px",
                    formatter: (cell, row) => {
                        return gridjs.html(`
                            <div class="d-flex gap-2">
                                <button class="btn btn-sm btn-info btn-edit-cs" data-id="${row.cells[0].data}"><i class="bx bx-edit"></i></button>
                                <button class="btn btn-sm btn-danger btn-delete-cs" data-id="${row.cells[0].data}"><i class="bx bx-trash"></i></button>
                            </div>
                        `);
                    }
                }
            ],
            server: {
                url: API_CS,
                then: data => data.data.map(item => [
                    item.id,
                    item.name,
                    item.color,
                    item.min_score,
                    item.max_score,
                    item.min_score,
                    item.max_score,
                    null
                ])
            },
            sort: true,
            pagination: { limit: 10 },
            style: {
                table: { 'white-space': 'nowrap' },
                th: { 'background-color': '#f8f9fa', 'text-align': 'center' },
                td: { 'vertical-align': 'middle', 'text-align': 'center' }
            }
        }).render(document.getElementById("table-confidence-score"));
    }

    // Helper to determine text color based on background
    function getContrastYIQ(hexcolor){
        hexcolor = hexcolor.replace("#", "");
        var r = parseInt(hexcolor.substr(0,2),16);
        var g = parseInt(hexcolor.substr(2,2),16);
        var b = parseInt(hexcolor.substr(4,2),16);
        var yiq = ((r*299)+(g*587)+(b*114))/1000;
        return (yiq >= 128) ? 'black' : 'white';
    }

    document.getElementById('btnAddCS').addEventListener('click', () => {
        formCS.reset();
        document.getElementById('csID').value = '';
        document.getElementById('modalCSLabel').innerText = 'Tambah Confidence Score';
        modalCS.show();
    });

    document.addEventListener('click', async (e) => {
        if (e.target.closest('.btn-edit-cs')) {
            const id = e.target.closest('.btn-edit-cs').dataset.id;
            editCS(id);
        }
        if (e.target.closest('.btn-delete-cs')) {
            const id = e.target.closest('.btn-delete-cs').dataset.id;
            deleteCS(id);
        }
    });

    async function editCS(id) {
        try {
            const res = await fetch(`${API_CS}`);
            const json = await res.json();
            const item = json.data.find(x => x.id == id);
            if (item) {
                document.getElementById('csID').value = item.id;
                document.getElementById('csName').value = item.name;
                document.getElementById('csColor').value = item.color;
                document.getElementById('csColorPicker').value = item.color;
                document.getElementById('csMin').value = item.min_score;
                document.getElementById('csMax').value = item.max_score;
                document.getElementById('modalCSLabel').innerText = 'Edit Confidence Score';
                modalCS.show();
            }
        } catch (err) { console.error(err); }
    }

    formCS.addEventListener('submit', async (e) => {
        e.preventDefault();
        const id = document.getElementById('csID').value;
        const data = {
            name: document.getElementById('csName').value,
            color: document.getElementById('csColor').value,
            min_score: parseInt(document.getElementById('csMin').value),
            max_score: parseInt(document.getElementById('csMax').value)
        };

        const method = id ? 'PUT' : 'POST';
        const url = id ? `${API_CS}/${id}` : API_CS;

        try {
            const res = await fetch(url, {
                method: method,
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
            const json = await res.json();
            if (json.status === 'success') {
                modalCS.hide();
                grid.forceRender();
                Swal.fire('Berhasil', 'Data berhasil disimpan', 'success');
            }
        } catch (err) { console.error(err); }
    });

    async function deleteCS(id) {
        const result = await Swal.fire({
            title: 'Apakah Anda yakin?',
            text: "Data akan dihapus permanen!",
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#d33',
            cancelButtonColor: '#3085d6',
            confirmButtonText: 'Ya, Hapus!'
        });

        if (result.isConfirmed) {
            try {
                const res = await fetch(`${API_CS}/${id}`, { method: 'DELETE' });
                const json = await res.json();
                if (json.status === 'success') {
                    grid.forceRender();
                    Swal.fire('Terhapus!', 'Data berhasil dihapus.', 'success');
                }
            } catch (err) { console.error(err); }
        }
    }

    initGrid();
});
