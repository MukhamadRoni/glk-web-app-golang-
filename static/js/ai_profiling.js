document.addEventListener('DOMContentLoaded', function () {
    const API_PROFILING = "/api/v1/ai/profiling";
    const modalElement = document.getElementById('modalProfiling');
    const modalProfiling = new bootstrap.Modal(modalElement);
    const formProfiling = document.getElementById('formProfiling');
    const fileInput = document.getElementById('skillFile');
    const btnSubmit = document.getElementById('btnSubmitSkill');

    let grid;

    function initGrid() {
        grid = new gridjs.Grid({
            columns: [
                { name: "ID", hidden: true },
                { name: "Nama Skill", width: "250px", formatter: (cell) => gridjs.html(`<strong>${cell || '-'}</strong>`) },
                { name: "Keterangan", width: "400px", formatter: (cell) => {
                    if (!cell) return "-";
                    return cell.length > 100 ? cell.substring(0, 100) + '...' : cell;
                }},
                { name: "Dokumen (.md)", width: "150px", formatter: (cell) => {
                    if (!cell) return "-";
                    return gridjs.html(`<a href="${cell}" target="_blank" class="btn btn-xs btn-outline-info"><i class="bx bx-link-external"></i> Link Drive</a>`);
                }},
                {
                    name: "Aksi",
                    width: "150px",
                    formatter: (cell, row) => {
                        return gridjs.html(`
                            <div class="d-flex gap-2">
                                <button class="btn btn-sm btn-info btn-edit-skill" data-id="${row.cells[0].data}"><i class="bx bx-edit"></i></button>
                                <button class="btn btn-sm btn-danger btn-delete-skill" data-id="${row.cells[0].data}"><i class="bx bx-trash"></i></button>
                            </div>
                        `);
                    }
                }
            ],
            server: {
                url: API_PROFILING,
                then: data => {
                    console.log("Profiling API Response:", data);
                    return data.data.map(item => [
                        item.id,
                        item.name,
                        item.keterangan,
                        item.url,
                        null
                    ]);
                }
            },
            sort: true,
            pagination: { limit: 10 },
            style: {
                table: { 'white-space': 'nowrap' },
                th: { 'background-color': '#f8f9fa', 'text-align': 'center' },
                td: { 'vertical-align': 'middle' }
            }
        }).render(document.getElementById("table-profiling"));
    }

    document.getElementById('btnAddSkill').addEventListener('click', () => {
        formProfiling.reset();
        document.getElementById('skillID').value = '';
        document.getElementById('urlDisplaySkill').style.display = 'none';
        document.getElementById('modalProfilingLabel').innerText = 'Tambah Profiling Skill Baru';
        modalProfiling.show();
    });

    document.addEventListener('click', async (e) => {
        if (e.target.closest('.btn-edit-skill')) {
            const id = e.target.closest('.btn-edit-skill').dataset.id;
            editSkill(id);
        }
        if (e.target.closest('.btn-delete-skill')) {
            const id = e.target.closest('.btn-delete-skill').dataset.id;
            deleteSkill(id);
        }
    });

    async function editSkill(id) {
        try {
            const res = await fetch(`${API_PROFILING}`);
            const json = await res.json();
            const item = json.data.find(x => x.id == id);
            if (item) {
                document.getElementById('skillID').value = item.id;
                document.getElementById('skillName').value = item.name;
                document.getElementById('skillKeterangan').value = item.keterangan;
                document.getElementById('skillURL').value = item.url;
                document.getElementById('btnViewSkillFile').href = item.url;
                document.getElementById('urlDisplaySkill').style.display = 'block';
                document.getElementById('modalProfilingLabel').innerText = 'Edit Profiling Skill';
                modalProfiling.show();
            }
        } catch (err) { console.error(err); }
    }

    formProfiling.addEventListener('submit', async (e) => {
        e.preventDefault();

        let fileUrl = document.getElementById('skillURL').value;

        if (fileInput.files.length > 0) {
            btnSubmit.disabled = true;
            btnSubmit.innerHTML = `<span class="spinner-border spinner-border-sm me-1"></span> Mengunggah File...`;

            try {
                fileUrl = await uploadToGDriveProfiling(fileInput.files[0]);
            } catch (err) {
                Swal.fire('Gagal Upload', err.message, 'error');
                btnSubmit.disabled = false;
                btnSubmit.innerHTML = `Simpan Data`;
                return;
            }
        }

        const id = document.getElementById('skillID').value;
        const data = {
            name: document.getElementById('skillName').value,
            keterangan: document.getElementById('skillKeterangan').value,
            url: fileUrl
        };

        const method = id ? 'PUT' : 'POST';
        const url = id ? `${API_PROFILING}/${id}` : API_PROFILING;

        try {
            const res = await fetch(url, {
                method: method,
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
            const json = await res.json();
            if (json.status === 'success') {
                modalProfiling.hide();
                grid.forceRender();
                Swal.fire('Berhasil', 'Data Profiling Skill berhasil disimpan', 'success');
            }
        } catch (err) { console.error(err); } finally {
            btnSubmit.disabled = false;
            btnSubmit.innerHTML = `Simpan Data`;
        }
    });

    async function uploadToGDriveProfiling(file) {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.readAsDataURL(file);
            reader.onload = async () => {
                const base64 = reader.result.split(',')[1];
                const payload = {
                    filename: `${document.getElementById('skillName').value.replace(/ /g, '_')}_${new Date().getTime()}.md`,
                    fileData: base64,
                    mimeType: 'text/markdown'
                };

                try {
                    const res = await fetch('/api/v1/ai/profiling/upload-proxy', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(payload)
                    });
                    const json = await res.json();
                    if (json.success) resolve(json.url);
                    else reject(new Error(json.message));
                } catch (err) { reject(err); }
            };
            reader.onerror = error => reject(error);
        });
    }

    async function deleteSkill(id) {
        const result = await Swal.fire({
            title: 'Hapus Data Profiling?',
            text: "Data skill ini akan dihapus permanen!",
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#d33',
            confirmButtonText: 'Ya, Hapus!'
        });

        if (result.isConfirmed) {
            try {
                const res = await fetch(`${API_PROFILING}/${id}`, { method: 'DELETE' });
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
