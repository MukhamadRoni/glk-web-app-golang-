(function () {
    const API_MAPEL = "/api/v1/mapel";
    const API_JP = "/api/v1/jenis-pendidikan";
    const tbodyMapel = document.querySelector("#tableMapel tbody");
    const selectJP = document.getElementById("mapelJenisPendidikan");

    // Fetch Jenis Pendidikan for select dropdown
    async function loadJP() {
        try {
            const res = await fetch(API_JP + "?_t=" + new Date().getTime(), { 
                credentials: "same-origin",
                cache: "no-store" 
            });
            const json = await res.json();
            
            selectJP.innerHTML = '<option value="">-- Pilih Jenis Pendidikan --</option>';
            if (json && json.data && Array.isArray(json.data)) {
                json.data.forEach(jp => {
                    const opt = document.createElement('option');
                    opt.value = jp.id;
                    opt.text = `${jp.name} (${jp.jenis_pendidikan})`;
                    selectJP.appendChild(opt);
                });
            }
        } catch (error) {
            console.error("Error loading JP:", error);
        }
    }

    // Fetch Mata Pelajaran
    async function loadMapel() {
        try {
            const res = await fetch(API_MAPEL + "?_t=" + new Date().getTime(), { 
                credentials: "same-origin",
                cache: "no-store"
            });
            const json = await res.json();
            
            tbodyMapel.innerHTML = "";

            if (json.data && json.data.length > 0) {
                json.data.forEach((m, index) => {
                    const isChecked = m.active === 'Y' || m.active === 'T' ? 'checked' : '';
                    const jpName = m.jenis_pendidikan_name ? m.jenis_pendidikan_name : '<span class="text-muted fst-italic">Belum Diatur</span>';
                    const jpId = m.jenis_pendidikan_id ? m.jenis_pendidikan_id : '';
                    const requirements = m.requirements ? m.requirements : '<span class="text-muted small italic">Tidak ada persyaratan khusus</span>';
                    
                    tbodyMapel.innerHTML += `
                        <tr>
                            <td>${index + 1}</td>
                            <td><span class="badge bg-soft-primary text-primary">${jpName}</span></td>
                            <td><code class="text-dark fw-bold">${m.code}</code></td>
                            <td><span class="fw-medium">${m.name}</span></td>
                            <td><small>${requirements}</small></td>
                            <td>
                                <div class="form-check form-switch form-switch-md" dir="ltr">
                                    <input type="checkbox" class="form-check-input switch-active-mapel" data-id="${m.id}" ${isChecked}>
                                    <label class="form-check-label small ms-1">${m.active === 'T' || m.active === 'Y' ? 'Aktif' : 'Non-aktif'}</label>
                                </div>
                            </td>
                            <td>
                                <div class="d-flex gap-2">
                                    <button class="btn btn-sm btn-soft-info waves-effect waves-light btn-edit-mapel" 
                                        data-id="${m.id}" 
                                        data-jp-id="${jpId}" 
                                        data-code="${m.code}" 
                                        data-name="${m.name}" 
                                        data-req="${m.requirements || ''}"
                                        title="Edit">
                                        <i class="bx bx-pencil"></i>
                                    </button>
                                    <button class="btn btn-sm btn-soft-danger waves-effect waves-light btn-delete-mapel" data-id="${m.id}" title="Hapus">
                                        <i class="bx bx-trash"></i>
                                    </button>
                                </div>
                            </td>
                        </tr>
                    `;
                });
            } else {
                tbodyMapel.innerHTML = `<tr><td colspan="7" class="text-center">Belum ada data Mata Pelajaran</td></tr>`;
            }
        } catch (error) {
            console.error("Error loading mapel:", error);
        }
    }

    // Toggle Active Status
    document.addEventListener("change", async function(e) {
        if (e.target.classList.contains("switch-active-mapel")) {
            const id = e.target.getAttribute("data-id");
            const newStatus = e.target.checked ? "T" : "F";
            
            try {
                Swal.fire({ title: 'Mohon Tunggu', html: 'Sedang memproses...', allowOutsideClick: false, didOpen: () => { Swal.showLoading() } });
                const res = await fetch(`${API_MAPEL}/${id}/active`, {
                    method: "PATCH",
                    headers: { "Content-Type": "application/json" },
                    credentials: "same-origin",
                    body: JSON.stringify({ active: newStatus })
                });
                
                if (!res.ok) {
                    const errData = await res.json();
                    Swal.fire('Gagal', "Gagal mengubah status: " + (errData.message || errData.error || "Unknown Error"), 'error');
                    e.target.checked = !e.target.checked; // Revert visually
                } else {
                    await loadMapel();
                    Swal.close();
                }
            } catch (error) {
                console.error(error);
                Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
                e.target.checked = !e.target.checked; // Revert visually
            }
        }
    });

    // Event Delegation for Create/Edit/Delete Buttons
    document.addEventListener("click", async function(e) {
        // Create Mapel
        const btnCreateMapel = e.target.closest('#btnCreateMapel');
        if (btnCreateMapel) {
            e.preventDefault();
            const jpId = document.getElementById("mapelJenisPendidikan").value;
            const code = document.getElementById("mapelCode").value;
            const name = document.getElementById("mapelName").value;
            const requirements = document.getElementById("mapelRequirements").value;

            if (!jpId || !code || !name) {
                Swal.fire('Peringatan', 'Jenis Pendidikan, Kode, dan Nama wajib diisi!', 'warning');
                return;
            }

            try {
                Swal.fire({ title: 'Mohon Tunggu', html: 'Sedang memproses...', allowOutsideClick: false, didOpen: () => { Swal.showLoading() } });
                const res = await fetch(API_MAPEL, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    credentials: "same-origin",
                    body: JSON.stringify({ 
                        jenis_pendidikan_id: parseInt(jpId), 
                        code: code, 
                        name: name,
                        requirements: requirements
                    })
                });
                if (res.ok) {
                    document.getElementById("mapelCode").value = "";
                    document.getElementById("mapelName").value = "";
                    document.getElementById("mapelRequirements").value = "";
                    await loadMapel();
                    Swal.fire('Berhasil', 'Berhasil menambahkan mata pelajaran!', 'success');
                } else {
                    const errData = await res.json();
                    Swal.fire('Gagal', "Gagal menambahkan data: " + (errData.message || errData.error || "Unknown Error"), 'error');
                }
            } catch (error) {
                console.error(error);
                Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
            }
        }

        // Delete Mapel
        const btnDeleteMapel = e.target.closest('.btn-delete-mapel');
        if (btnDeleteMapel) {
            const id = btnDeleteMapel.getAttribute('data-id');
            const confirm = await Swal.fire({
                title: 'Yakin ingin menghapus?',
                text: "Data mata pelajaran akan dihapus secara logika (soft-delete).",
                icon: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#d33',
                cancelButtonColor: '#3085d6',
                confirmButtonText: 'Ya, hapus!',
                cancelButtonText: 'Batal'
            });
            if (!confirm.isConfirmed) return;

            try {
                Swal.fire({ title: 'Mohon Tunggu', html: 'Sedang memproses...', allowOutsideClick: false, didOpen: () => { Swal.showLoading() } });
                const res = await fetch(`${API_MAPEL}/${id}`, { 
                    method: "DELETE",
                    credentials: "same-origin"
                });
                if (res.ok) {
                    await loadMapel();
                    Swal.fire('Terhapus!', 'Data berhasil dihapus.', 'success');
                } else {
                    Swal.fire('Gagal', 'Gagal menghapus data', 'error');
                }
            } catch (error) {
                console.error(error);
                Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
            }
        }

        // Edit Mapel
        const btnEditMapel = e.target.closest('.btn-edit-mapel');
        if (btnEditMapel) {
            const id = btnEditMapel.getAttribute('data-id');
            const currentJpId = btnEditMapel.getAttribute('data-jp-id');
            const currentCode = btnEditMapel.getAttribute('data-code');
            const currentName = btnEditMapel.getAttribute('data-name');
            const currentReq = btnEditMapel.getAttribute('data-req');
            
            const jpOptions = Array.from(selectJP.options)
                .filter(opt => opt.value !== "")
                .map(opt => `<option value="${opt.value}" ${opt.value === currentJpId ? 'selected' : ''}>${opt.text}</option>`)
                .join('');

            const { value: formValues } = await Swal.fire({
                title: 'Edit Mata Pelajaran',
                html:
                    `<div class="text-start">` +
                    `<label class="form-label small fw-bold">Jenis Pendidikan</label>` +
                    `<select id="swal-mapelJp" class="form-select mb-3">${jpOptions}</select>` +
                    `<label class="form-label small fw-bold">Kode</label>` +
                    `<input id="swal-mapelCode" class="form-control mb-3" placeholder="Kode" value="${currentCode}">` +
                    `<label class="form-label small fw-bold">Nama Mata Pelajaran</label>` +
                    `<input id="swal-mapelName" class="form-control mb-3" placeholder="Nama" value="${currentName}">` +
                    `<label class="form-label small fw-bold">Requirements</label>` +
                    `<textarea id="swal-mapelReq" class="form-control mb-3" placeholder="Requirements" rows="3">${currentReq}</textarea>` +
                    `</div>`,
                focusConfirm: false,
                showCancelButton: true,
                preConfirm: () => {
                    const jpId = document.getElementById('swal-mapelJp').value;
                    const code = document.getElementById('swal-mapelCode').value;
                    const name = document.getElementById('swal-mapelName').value;
                    const req = document.getElementById('swal-mapelReq').value;
                    if (!jpId || !code || !name) {
                        Swal.showValidationMessage('Jenis Pendidikan, Kode, dan Nama wajib diisi!');
                    }
                    return { jpId, code, name, req };
                }
            });

            if (formValues) {
                try {
                    Swal.fire({ title: 'Mohon Tunggu', html: 'Sedang memproses...', allowOutsideClick: false, didOpen: () => { Swal.showLoading() } });
                    const res = await fetch(`${API_MAPEL}/${id}`, {
                        method: "PUT",
                        headers: { "Content-Type": "application/json" },
                        credentials: "same-origin",
                        body: JSON.stringify({ 
                            jenis_pendidikan_id: parseInt(formValues.jpId), 
                            code: formValues.code, 
                            name: formValues.name,
                            requirements: formValues.req
                        })
                    });
                    if (res.ok) {
                        await loadMapel();
                        Swal.fire('Tersimpan!', 'Data berhasil diubah.', 'success');
                    } else {
                        const errData = await res.json();
                        Swal.fire('Gagal', 'Gagal mengubah data: ' + (errData.message || errData.error || "Unknown Error"), 'error');
                    }
                } catch (error) {
                    console.error(error);
                    Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
                }
            }
        }
    });

    // Initial load
    loadJP().then(() => loadMapel());
})();
