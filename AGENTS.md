# AGENTS.md — HOSIM-GO Engineering Rules

Dokumen ini berlaku untuk seluruh repository HOSIM-GO. Semua coding agent WAJIB membacanya sebelum membuat atau mengubah kode. Jika ada `AGENTS.md` yang lebih dekat dengan file yang dikerjakan, aturan yang lebih spesifik tersebut melengkapi dokumen ini; jika bertentangan, ikuti aturan yang paling dekat selama tidak merusak prinsip keselamatan, integritas data medis, dan batas modul di dokumen root ini.

Kata **WAJIB**, **DILARANG**, dan **BOLEH** bersifat normatif.

## 1. Tujuan arsitektur

HOSIM-GO adalah **pragmatic modular monolith** dengan **pragmatic Clean Architecture**:

- Satu deployment dan satu database PostgreSQL, tetapi kode dibagi menjadi modul bisnis dengan ownership yang tegas.
- Dependency mengarah ke business rules. Framework, HTTP, GORM, PostgreSQL, dan API eksternal adalah detail di sisi luar.
- Gunakan abstraksi hanya pada boundary yang nyata. Jangan membuat layer, interface, atau mapping tanpa manfaat konkret.
- Integritas data medis, jejak audit, otorisasi kontekstual, dan idempotensi lebih penting daripada mempersingkat kode.
- Jangan memaksakan satu pola untuk semua fitur. CRUD master sederhana dan transaksi bisnis memakai alur yang berbeda.

Arsitektur ini bukan izin untuk membuat “distributed monolith” di dalam satu repository. Modul harus berkomunikasi melalui contract yang sempit, bukan membaca tabel atau memakai repository internal modul lain secara sembarang.

## 2. Pilih jalur implementasi sebelum menulis kode

Agent WAJIB mengklasifikasikan perubahan sebagai salah satu dari dua jalur berikut.

### Jalur A — simple master data

Gunakan:

```text
HTTP Handler -> Service -> Repository -> PostgreSQL
             <- DTO/result <-
```

Jalur ini hanya BOLEH dipakai jika **semua** kondisi berikut benar:

- Fitur utamanya CRUD atau pencarian data referensi/konfigurasi.
- Operasi berpusat pada satu aggregate/tabel milik satu modul.
- Tidak mengorkestrasi modul lain.
- Tidak mempunyai state transition atau workflow bisnis yang berarti.
- Tidak memerlukan beberapa write yang harus atomik.
- Tidak memicu integrasi eksternal, outbox event, atau side effect asinkron.
- Aturan bisnisnya sederhana: validasi field, uniqueness, status aktif, dan referential check ringan.

Contoh yang biasanya cocok: department, service unit, tariff class, katalog referensi, atau master konfigurasi lain. “Master data” bukan pengecualian otomatis: perubahan bed availability, kuota, jadwal, tarif efektif, atau konfigurasi yang memicu workflow dapat memerlukan Jalur B.

Pada Jalur A, `Service` boleh mengoordinasikan validasi dan satu repository. Jika service mulai memegang banyak repository, transaction manager, state machine, audit klinis, atau client eksternal, hentikan dan refactor ke use case.

### Jalur B — business transaction atau cross-module workflow

Gunakan:

```text
HTTP Handler -> Application Use Case -> Domain Entity/Domain Service
                                      -> Repository ports
                                      -> Outbox/Audit ports
             <- Result <-
```

Pilih Jalur B jika **salah satu** kondisi berikut benar:

- Nama operasinya adalah tindakan/tujuan bisnis, misalnya `AdmitPatient`, `TransferPatient`, `DischargePatient`, `RecordVitalSigns`, `PrescribeMedication`, atau `DispenseMedication`.
- Operasi mengubah lifecycle/status dan mempunyai allowed transition.
- Operasi menyentuh beberapa aggregate, repository, atau modul.
- Beberapa write harus commit atau rollback bersama.
- Ada aturan klinis/finansial yang harus selalu benar.
- Ada audit trail bernilai medis/finansial.
- Ada integration event, transactional outbox, retry, atau idempotency requirement.
- Hasil bergantung pada actor, organisasi, fasilitas, unit layanan, encounter, atau assignment klinis.

Jangan menamai orchestration bisnis kompleks sebagai `FooService` hanya agar mengikuti pola lama. Satu file use case harus mewakili satu intent yang jelas, misalnya `admit_patient.go`, bukan `inpatient_service.go` yang berisi semua operasi.

### Read/query path

Query read-only BOLEH memakai query service/repository yang lebih langsung bila tidak melewati otorisasi dan module ownership. Jangan membangun aggregate lengkap hanya untuk list/report. Namun query tetap WAJIB scoped, tidak boleh membocorkan data pasien, dan tidak boleh menjadikan handler sebagai tempat SQL atau business rule.

## 3. Dependency direction

Dependency yang diizinkan:

```text
transport/http (Gin, request/response DTO)
                 |
                 v
application (use case, command/result, ports)
                 |
                 v
domain (entity, value object, invariant, domain service)

infrastructure/persistence (GORM/PostgreSQL) --implements--> repository ports
infrastructure/integration (HTTP clients)    --implements--> integration ports
cmd/api                                      --wires-------> semua implementasi
```

Aturan tegas:

- Domain DILARANG mengimpor Gin, GORM, driver database, HTTP client, atau package transport.
- Application DILARANG bergantung pada concrete GORM repository atau concrete external client.
- Handler hanya bergantung pada service/use case contract yang diperlukan.
- Infrastructure boleh bergantung ke domain/application untuk mengimplementasikan port, bukan sebaliknya.
- Wiring dilakukan di composition root (`cmd/api` atau package bootstrap/server yang sudah dipakai repository).
- Hindari global database/client/service locator. Inject dependency melalui constructor.
- Semua operasi I/O menerima `context.Context` sebagai argumen pertama dan meneruskannya ke dependency.
- Jangan melakukan refactor folder besar hanya demi menyerupai diagram ini. Pertahankan layout existing selama dependency direction dan tanggung jawabnya benar.

## 4. Bentuk data dan mapping antar-layer

Gunakan tipe yang berbeda untuk tanggung jawab yang berbeda:

| Tipe | Milik layer | Isi dan aturan |
|---|---|---|
| Request DTO | HTTP/transport | Bentuk payload masuk, `json`/binding tags, validasi sintaks dan format. Jangan berisi GORM tags. |
| Command/Query | Application | Input use case yang sudah dinormalisasi, ditambah actor/scope/idempotency context yang diperlukan. Tidak bergantung Gin. |
| Domain entity/value object | Domain | State dan perilaku bisnis; menjaga invariant. Bukan kontrak HTTP. |
| Persistence model/record | Infrastructure | GORM tags, nama kolom, relasi persistence, nullable types, timestamps. |
| Result | Application | Hasil use case yang stabil bagi caller; tidak membawa `*gorm.DB` atau model internal. |
| Response DTO | HTTP/transport | Kontrak API keluar; mapping eksplisit dari result. Jangan expose field internal/sensitif. |

Alur write yang benar:

```text
JSON request
  -> Request DTO (bind + syntactic validation)
  -> Command (normalization + authenticated actor/scope)
  -> Use case/service
  -> Domain method/domain service (business invariant)
  -> Repository port
  -> GORM persistence record
  -> PostgreSQL

PostgreSQL
  -> persistence record
  -> domain/result
  -> Response DTO
  -> JSON response
```

DILARANG bind request langsung ke entity atau GORM model. DILARANG mengembalikan entity/GORM model langsung sebagai response. Field seperti password hash, internal note, deleted timestamp, raw integration payload, dan metadata audit tidak boleh bocor karena serialisasi otomatis.

Untuk CRUD yang benar-benar sederhana, entity dan persistence model BOLEH berupa tipe yang sama jika pemisahan hanya menghasilkan mapping mekanis tanpa melindungi boundary apa pun. Walaupun demikian, request/response DTO tetap terpisah dan GORM concerns tidak boleh masuk ke handler/service contract. Begitu persistence shape mulai mendikte domain atau data medis sensitif rawan terekspos, pisahkan record dari entity.

## 5. Handler, service, dan use case

### Handler

Handler hanya boleh:

1. Membaca path/query/header/body.
2. Bind dan memvalidasi sintaks/format request.
3. Mengambil authenticated principal, correlation ID, dan request metadata dari context.
4. Membentuk command/query.
5. Memanggil tepat satu application entry point untuk satu intent.
6. Menerjemahkan typed error ke status dan error response standar.
7. Menerjemahkan result ke response DTO.

Handler DILARANG berisi transaction, query GORM, state transition, keputusan klinis, cross-module orchestration, atau external HTTP call.

### Service pada Jalur A

Service mengelola aturan CRUD sederhana, normalisasi, uniqueness check bila diperlukan, dan pemanggilan repository. Jangan membuat interface service kecuali ada caller/boundary yang benar-benar membutuhkan substitution; interface sebaiknya didefinisikan oleh consumer, bukan otomatis untuk setiap struct.

### Use case pada Jalur B

Use case adalah application orchestrator untuk satu tujuan bisnis. Use case WAJIB:

- Memvalidasi actor dan contextual authorization sebelum membaca/mengubah data sensitif.
- Memuat state yang dibutuhkan melalui port.
- Memanggil perilaku entity/domain service untuk aturan bisnis.
- Menentukan transaction boundary.
- Menyimpan semua perubahan atomik.
- Menulis audit record dan outbox event dalam transaksi yang sama bila diwajibkan.
- Mengembalikan result yang tidak bergantung pada transport atau persistence.

Contoh bentuk, bukan template yang harus disalin mentah:

```go
type AdmitPatientCommand struct {
	PatientID      string
	ServiceUnitID  string
	BedID          string
	DoctorID       string
	Actor          ActorContext
	IdempotencyKey string
}

type AdmitPatient struct {
	tx         TransactionManager
	patients   PatientReader
	encounters EncounterRepository
	beds       BedRepository
	admissions AdmissionRepository
	audit      AuditWriter
	outbox     OutboxWriter
}

func (uc *AdmitPatient) Execute(ctx context.Context, cmd AdmitPatientCommand) (AdmitPatientResult, error) {
	// authorize -> load -> apply domain rules -> atomic save/audit/outbox -> result
}
```

## 6. Transaction boundary dan concurrency

- Use case adalah pemilik transaction boundary untuk business transaction. Handler dan entity DILARANG membuka transaksi.
- Repository tunggal tidak boleh diam-diam memulai transaction untuk workflow multi-repository.
- Gunakan `TransactionManager`/unit-of-work abstraction atau mekanisme project yang setara agar application tidak bergantung langsung pada GORM.
- Semua repository yang dipanggil di dalam transaksi WAJIB memakai handle transaksi yang sama.
- Transaksi harus sesingkat mungkin: authorize/load state yang relevan, lock bila perlu, validate invariant, write, audit/outbox, commit.
- Jangan melakukan network call, file I/O lambat, menunggu message broker, atau pekerjaan CPU berat di dalam DB transaction.
- Untuk resource yang diperebutkan—bed, stock, nomor dokumen, payment state—tangani race dengan database constraint, atomic conditional update, optimistic version, atau row lock (`FOR UPDATE`) yang disengaja. Check-then-write tanpa proteksi concurrency DILARANG.
- Gunakan isolation/locking yang sesuai kasus; jangan menaikkan isolation global tanpa bukti.
- Jika request dapat diulang, definisikan idempotency key dan simpan hasil/status secara atomik. Retry tidak boleh menggandakan admission, dispensing, billing, audit utama, atau outbox event.

## 7. Domain invariants

Aturan yang harus selalu benar ditempatkan pada entity/value object atau domain service, bukan hanya handler/UI. Contoh:

- Encounter yang sudah selesai tidak boleh menerima clinical record baru kecuali ada workflow koreksi yang sah.
- Bed tidak boleh ditempati dua admission aktif.
- Transfer hanya boleh dari admission aktif dan tujuan harus tersedia.
- Dispensing tidak boleh melebihi kuantitas yang diizinkan tanpa aturan override yang eksplisit.
- State transition harus berasal dari status yang diizinkan.

Gunakan entity method seperti `admission.TransferTo(...)` atau domain service bila aturan melibatkan beberapa entity dan tidak alami dimiliki satu entity. Jangan membuat domain service menjadi kumpulan helper stateless tanpa business meaning.

Validasi berlapis:

- Transport: required, format, enum parsing, batas ukuran.
- Application: existence, authorization, orchestration precondition.
- Domain: invariant dan state transition.
- Database: `NOT NULL`, `UNIQUE`, `CHECK`, foreign key, dan index sebagai pertahanan terakhir.

## 8. Repository contracts

- Repository interface/port dimiliki oleh domain/application yang menggunakannya, bukan oleh package infrastruktur generik.
- Method harus memakai bahasa bisnis: `FindActiveByPatient`, `ReserveAvailableBed`, `Save`, bukan sekadar repository generik dengan semua operasi untuk semua entity.
- Setiap method I/O menerima `context.Context`.
- Port tidak boleh mengekspos `*gorm.DB`, GORM clauses, SQL rows, atau persistence record.
- Bedakan “tidak ditemukan” dari infrastructure failure dengan typed/sentinel error.
- List query harus punya pagination dan batas maksimum. Filter/sort yang diizinkan harus eksplisit; jangan meneruskan nama kolom mentah dari request.
- Cross-module consumer tidak boleh mengambil concrete repository modul lain. Pakai exported query/use-case contract atau narrow port yang disediakan owner.
- Hindari `BaseRepository[T]` yang menghapus business semantics dan membuat semua tabel dapat diakses semua modul.

## 9. GORM dan PostgreSQL

- PostgreSQL adalah source of truth produksi. SQLite BOLEH untuk test cepat yang tidak bergantung perilaku database, tetapi tidak cukup untuk membuktikan query, constraint, locking, JSONB, UUID, timestamp, index, atau transaction semantics PostgreSQL.
- Semua query GORM harus memakai context (`db.WithContext(ctx)` atau handle ekuivalen).
- Pilih kolom secara sadar untuk update. Hindari `Save`/mass update yang tanpa sengaja menimpa zero value atau field yang tidak boleh diubah.
- Jangan memakai map/request body langsung sebagai update payload.
- Tangani `gorm.ErrRecordNotFound`, unique/FK/check violation, timeout, dan cancellation secara konsisten lalu terjemahkan ke application error.
- Cegah N+1. Gunakan preload/join/batch secara terukur, dan hanya load relasi yang dibutuhkan.
- Gunakan pointer atau nullable type secara sengaja; bedakan nilai kosong dari nilai tidak tersedia.
- Semua waktu bisnis disimpan dengan timezone yang jelas; utamakan UTC di persistence dan konversi di boundary. Jangan bergantung pada timezone proses/database secara implisit.
- Soft delete bukan default. Pakai hanya jika semantics bisnis membutuhkannya, lalu pastikan unique constraint, query scope, audit, dan restore behavior benar.
- Nama tabel/kolom/index/constraint harus eksplisit dan stabil bila convention GORM berisiko berubah.
- `AutoMigrate` DILARANG sebagai mekanisme perubahan schema saat aplikasi start, terutama di staging/production.
- Jangan mengandalkan hook GORM tersembunyi untuk workflow penting. Business side effect harus terlihat di use case/domain; hook hanya untuk concern persistence lokal yang aman dan teruji.

## 10. Goose migrations

Semua perubahan schema WAJIB melalui file Goose di `migrations/`.

- Buat migration baru; DILARANG mengedit migration yang sudah pernah diterapkan/shared hanya agar environment lokal lolos.
- Gunakan urutan/naming yang mengikuti convention repository.
- Sertakan `-- +goose Up` dan `-- +goose Down` bila rollback aman. Jika down migration berisiko kehilangan data, dokumentasikan keputusan dan sediakan strategi recovery yang realistis daripada rollback palsu.
- Definisikan FK, `NOT NULL`, `UNIQUE`, `CHECK`, index, dan delete/update behavior secara eksplisit.
- Untuk tabel besar atau perubahan berisiko, gunakan pola expand -> deploy/backfill -> switch -> contract. Hindari table rewrite/long lock pada deployment tunggal.
- Backfill harus deterministic, restartable/idempotent bila mungkin, dan tidak diam-diam mengubah data klinis tanpa audit/rencana verifikasi.
- Seed data referensi dipisahkan dari schema migration kecuali data tersebut mutlak diperlukan oleh schema/constraint.
- Verifikasi migration dari database kosong dan upgrade dari schema sebelumnya. Pastikan query lama/baru kompatibel selama rollout yang direncanakan.

## 11. Module ownership dan dependency rules

Modul utama dan ownership konseptual:

- `organization`: department, service unit, room, bed, storage, depo, referral facility.
- `patient`: identitas dan demographic pasien; bukan pemilik seluruh rekam medis.
- `practitioner`: identitas/credential tenaga kesehatan.
- `catalog`: kode dan referensi medis.
- `finance`: payer/customer, tariff class, tariff policy/reference.
- `clinical`: shared EHR capabilities dan clinical facts.
- `emergency`, `outpatient`, `inpatient`: workflow pelayanan masing-masing.
- `pharmacy`: prescription/dispensing workflow sesuai boundary yang disepakati.
- `inventory`: stock, movement, reservation, valuation/logistics.
- `billing`: charge, invoice, payment/cashier workflow.
- `integration`: adapter SATUSEHAT, BPJS, LIS, PACS dan delivery infrastructure.
- `audit`: append-only audit capability.
- `auth`: identity/session/role-permission capability.

Aturan antar-modul:

- Satu modul adalah satu owner untuk mutasi aggregate/tabelnya.
- Modul lain DILARANG update tabel owner secara langsung, memakai GORM model internalnya, atau mengimpor package `persistence`-nya.
- Cross-module write dilakukan melalui application contract/use case owner atau event yang didefinisikan jelas.
- Cross-module read memakai exported query contract/read model yang sempit. Direct join lintas tabel hanya BOLEH untuk read-only reporting/query terkontrol, tetap dimiliki package query yang jelas dan tidak dipakai untuk mutasi.
- DILARANG membuat circular dependency. Bila A dan B saling membutuhkan, ekstrak contract/value yang stabil atau pindahkan orchestration ke module/use case pemilik workflow; jangan membuat package `common` sebagai tempat membuang semua tipe.
- `pkg/` hanya untuk library teknis yang benar-benar generik dan stabil. Business concepts tetap di `internal/<module>`.

### Clinical capability vs care-setting workflow

`clinical` memiliki capability/fact klinis yang reusable lintas setting, misalnya Encounter, Observation, Condition, Procedure, Allergy, clinical note, atau order representation yang memang shared.

`emergency`, `outpatient`, dan `inpatient` memiliki alur pelayanan, state transition, dan aturan setting masing-masing:

- Emergency: triage, arrival/acuteness, emergency disposition.
- Outpatient: registration/queue/visit flow, consultation completion.
- Inpatient: admission, bed assignment, transfer, discharge.

DILARANG menaruh semua workflow tersebut di `clinical` hanya karena menghasilkan data klinis. Sebaliknya, jangan menduplikasi Observation/Condition/Procedure untuk tiap setting jika semantics-nya sama. Workflow module mengorkestrasi shared clinical capability melalui contract-nya.

### Encounter sebagai central service context

Encounter adalah konteks pusat satu episode/interaksi pelayanan, bukan tabel serbaguna dan bukan pengganti admission/visit/triage:

- Clinical records yang dibuat dalam pelayanan WAJIB mereferensikan `EncounterID` bila secara bisnis memang berada dalam encounter.
- Encounter menghubungkan patient, care setting/type, organization/service unit, responsible participants, period, dan status.
- Emergency visit, outpatient visit, dan inpatient admission memiliki lifecycle sendiri tetapi terkait ke encounter yang sesuai.
- Hanya owner contract yang boleh membuat/mengubah/menutup Encounter. Workflow module meminta perubahan melalui contract tersebut, bukan update tabel `encounters` secara langsung.
- Status Encounter dan status workflow harus disinkronkan dalam use case atomik atau lewat event/outbox yang idempotent; definisikan source of truth untuk setiap status.
- Jangan memakai Encounter sebagai tempat menumpuk field khusus IGD/rawat jalan/rawat inap. Field khusus tetap milik workflow module.

## 12. Integration adapters: SATUSEHAT, BPJS, LIS, PACS

Application/domain mendefinisikan port berdasarkan kebutuhan bisnis. `integration/<provider>` mengimplementasikan adapter dan mapping ke protokol provider.

- Domain DILARANG mengenal FHIR HTTP, BPJS endpoint, LIS/PACS vendor schema, OAuth token response, atau SDK provider.
- Pisahkan canonical/internal model dari provider DTO. Mapping provider harus eksplisit dan teruji.
- Semua call memakai timeout, context cancellation, correlation ID, structured logging yang aman, dan error classification (retryable vs permanent).
- Terapkan idempotency/deduplication untuk request dan callback/webhook.
- Validasi authenticity callback, signature bila tersedia, replay protection, dan ownership patient/encounter sebelum memproses payload.
- Jangan log token, credential, NIK penuh, detail klinis, atau raw payload sensitif. Bila raw payload wajib disimpan untuk traceability, enkripsi/batasi akses dan tetapkan retention.
- Retry memakai bounded exponential backoff dan dead-letter/manual reconciliation path; jangan retry validation/permanent error tanpa batas.
- Simpan external identifier dan sync status dengan constraint yang mencegah duplikasi.

### Transactional outbox

External HTTP call DILARANG dilakukan di dalam DB transaction.

Untuk efek eksternal akibat perubahan state:

1. Use case mengubah state internal.
2. Use case menyimpan outbox message dalam transaksi database yang sama.
3. Commit.
4. Worker mengambil outbox message.
5. Adapter memanggil provider.
6. Worker mencatat sukses/gagal, retry schedule, dan external reference secara idempotent.

Outbox event minimal memiliki ID unik, event type/version, aggregate ID, occurred time, correlation/causation ID, payload minimal, attempt/status, dan next-attempt metadata. Consumer harus idempotent karena delivery bersifat at-least-once. Jangan memasukkan data medis berlebih ke payload hanya demi kemudahan.

## 13. Audit trail

Audit berbeda dari application log. Log membantu operasi; audit membuktikan siapa melakukan apa pada data sensitif.

Untuk create/update/delete, state transition, access sensitif, override, export, dan integration submission yang relevan, catat setidaknya:

- actor/user/practitioner dan role efektif;
- organization/facility/service unit scope;
- patient ID dan encounter ID bila relevan;
- action, resource type, resource ID, serta outcome;
- timestamp, request/correlation ID, source/channel;
- reason/justification untuk override atau break-glass;
- perubahan before/after atau field-level diff yang aman dan proporsional.

Aturan:

- Audit bernilai kritis ditulis append-only dalam transaction yang sama dengan perubahan state, atau melalui mekanisme durable yang memberi jaminan setara.
- Audit failure pada operasi yang diwajibkan harus menggagalkan transaksi; jangan diam-diam lanjut.
- Audit record tidak boleh diubah/dihapus lewat CRUD biasa.
- Jangan menyimpan secret, password, token, atau payload medis besar secara buta di audit.
- Read access ke rekam medis sensitif juga dapat wajib diaudit; ikuti policy feature terkait.

## 14. Authentication, authorization, dan contextual access

- Middleware melakukan authentication dan parsing credential, tetapi use case adalah enforcement point untuk business authorization.
- Jangan mempercayai `actor_id`, role, organization, service unit, atau patient scope dari request body. Ambil dari authenticated context dan data server-side.
- Otorisasi harus mempertimbangkan action + resource + context, bukan role saja: organisasi/fasilitas aktif, unit layanan, assignment/relationship ke encounter, status encounter, dan purpose of use bila relevan.
- Default adalah deny. Endpoint baru tidak boleh aktif tanpa permission mapping yang jelas.
- Repository/query untuk data tenant/organisasi WAJIB menerima scope atau predicate yang membuat query tanpa scope sulit dilakukan. Hindari `FindByID(id)` global untuk resource sensitif bila ownership scope diperlukan.
- Cegah IDOR: setelah menemukan resource, tetap verifikasi resource berada dalam scope actor.
- Break-glass access harus eksplisit, time-bound bila policy mengharuskan, meminta alasan, dan menghasilkan audit prioritas tinggi.
- UI hiding bukan authorization. Semua aturan tetap ditegakkan backend.

## 15. Error handling dan API behavior

Definisikan error application/domain yang dapat diperiksa dengan `errors.Is`/`errors.As`, misalnya kategori:

- validation/business rule;
- unauthenticated;
- forbidden;
- not found;
- conflict/concurrent modification/idempotency conflict;
- dependency unavailable/timeout;
- internal.

Handler memetakan kategori tersebut secara konsisten ke response envelope dan HTTP status. Jangan menentukan status HTTP di domain/repository.

- Bungkus error dengan konteks menggunakan `%w`; jangan kehilangan cause.
- Jangan membandingkan string error.
- Jangan mengembalikan SQL, nama tabel, stack trace, provider credential, atau detail internal ke client.
- Log error sekali di boundary yang memiliki request/correlation context. Hindari logging berulang di setiap layer.
- Validation response harus stabil dan actionable tanpa membocorkan data.
- Hormati `context.Canceled` dan deadline; jangan mengubah cancellation menjadi 500 generik.
- Panic hanya untuk invariant bootstrap/configuration yang membuat proses tidak dapat berjalan, bukan untuk request error normal.

## 16. Struktur folder yang disarankan

Ikuti struktur existing dan lakukan perubahan bertahap. Contoh simple master data:

```text
internal/organization/department/
├── entity.go
├── repository.go          # interface/contract
├── service.go             # simple CRUD rules
├── handler.go             # Gin + HTTP DTO mapping
├── repository_gorm.go     # GORM implementation
└── *_test.go
```

Contoh workflow kompleks; subfolder boleh disesuaikan dengan convention yang sudah ada:

```text
internal/inpatient/admission/
├── domain/
│   ├── admission.go
│   ├── rules.go
│   └── repository.go
├── application/
│   ├── admit_patient.go
│   ├── transfer_patient.go
│   └── discharge_patient.go
├── transport/http/
│   ├── handler.go
│   └── dto.go
├── persistence/
│   ├── model.go
│   └── repository_gorm.go
└── *_test.go
```

Jika layout repository saat ini flat, struktur berikut juga valid dan lebih pragmatis:

```text
internal/inpatient/admission/
├── entity.go
├── repository.go
├── usecase/
│   ├── admit_patient.go
│   ├── transfer_patient.go
│   └── discharge_patient.go
├── handler.go
├── repository_gorm.go
└── *_test.go
```

Jangan membuat package satu file secara berlebihan. Pisahkan package ketika boundary/dependency memang perlu ditegakkan, bukan untuk mengejar bentuk diagram.

## 17. Testing expectations

Setiap perubahan behavior WAJIB disertai test pada level paling rendah yang dapat membuktikannya:

- Domain unit test: invariant, allowed/forbidden transition, boundary value; tanpa database/network.
- Service/use-case unit test: orchestration, authorization, transaction outcome, audit/outbox invocation, error propagation; gunakan fake/mock port yang fokus.
- Repository integration test: query, mapping, constraints, transaction, locking/concurrency terhadap PostgreSQL untuk behavior PostgreSQL-specific.
- Handler test: binding, authenticated context, error-to-status mapping, response contract; jangan mengulang seluruh domain test.
- Adapter contract test: request/response mapping, timeout, retry classification, idempotency, callback validation memakai fake server.
- Migration test: fresh up dan upgrade path; down hanya bila dijanjikan aman.
- Outbox worker test: duplicate delivery, retry, crash-after-send scenario, poison message/dead-letter behavior.

Minimal sebelum selesai:

```text
gofmt pada file Go yang berubah
go test ./...
go vet ./...
```

Jalankan test yang paling spesifik selama iterasi, lalu suite relevan/full suite sebelum handoff bila environment memungkinkan. Jika test membutuhkan PostgreSQL atau credential eksternal yang tidak tersedia, agent WAJIB menyebutkan secara eksplisit apa yang tidak dijalankan dan tetap menjalankan test lain yang tersedia. Jangan mengklaim test lulus jika tidak dijalankan.

Untuk bug fix, tambahkan regression test yang gagal sebelum fix dan lulus sesudah fix. Untuk jalur kritis, test juga forbidden path dan failure path, bukan hanya happy path.

## 18. Anti-patterns yang dilarang

- Fat handler atau `service.go` raksasa yang menjadi tempat semua orchestration.
- CRUD mental model untuk admission, transfer, discharge, prescribing, dispensing, billing, atau workflow klinis lain.
- Request DTO = domain entity = GORM model = response DTO.
- Domain mengimpor Gin/GORM atau mengembalikan HTTP status.
- Handler/repository memulai transaksi cross-module secara tersembunyi.
- External HTTP call di dalam DB transaction.
- Update tabel modul lain secara langsung.
- Generic repository/service yang menghapus bahasa dan invariant bisnis.
- Business rule hanya di UI/handler/database hook.
- `AutoMigrate` saat startup sebagai pengganti Goose.
- Mengubah migration yang sudah diterapkan.
- Check-then-write untuk bed/stock/payment tanpa concurrency guard.
- Query data pasien tanpa organization/context scope.
- Mempercayai actor/role/scope dari body request.
- Menelan error, `panic` untuk error request biasa, string matching error, atau selalu mengembalikan 500.
- Log/audit yang memuat token, password, NIK penuh, atau payload medis tanpa redaction dan kebutuhan jelas.
- Event/outbox tanpa idempotency dan versioning.
- Import cycle yang “diselesaikan” dengan memindahkan seluruh business type ke `pkg/common`.
- Abstraksi spekulatif, interface satu implementasi tanpa boundary/testing need, atau factory berlapis yang tidak menambah safety.

## 19. Checklist wajib bagi AI coding agent

### Sebelum mengubah kode

- [ ] Baca `AGENTS.md` root dan file instruksi terdekat.
- [ ] Inspeksi kode, test, migration, route, wiring, dan convention modul terkait; jangan menebak.
- [ ] Nyatakan secara internal: ini Jalur A atau Jalur B? Pastikan kriterianya.
- [ ] Tentukan owner aggregate/tabel dan modul yang hanya menjadi consumer.
- [ ] Identifikasi invariant, state transition, authorization scope, data sensitif, dan concurrency risk.
- [ ] Identifikasi apakah butuh transaction, audit, outbox, idempotency, atau migration.
- [ ] Cari contract/existing helper yang relevan sebelum membuat abstraksi baru.

### Saat mengimplementasikan

- [ ] Jaga dependency direction dan hindari import concrete infrastructure dari application/domain.
- [ ] Pisahkan HTTP DTO, command/result, domain type, dan persistence record sesuai kebutuhan boundary.
- [ ] Pastikan handler tipis dan satu intent memanggil satu application entry point.
- [ ] Letakkan invariant di domain dan transaction boundary di use case.
- [ ] Gunakan context, scope authorization, typed error, dan explicit field mapping.
- [ ] Untuk workflow lintas modul, gunakan contract owner; jangan query/update tabel internalnya.
- [ ] Tulis audit/outbox secara atomik bila diperlukan; jangan call provider di transaction.
- [ ] Tambahkan Goose migration untuk perubahan schema; jangan memakai AutoMigrate.
- [ ] Tambahkan test untuk happy path, forbidden/business failure, dan rollback/concurrency/idempotency bila relevan.

### Sebelum menyatakan selesai

- [ ] Review diff untuk accidental data exposure, mass assignment, unscoped query, N+1, dan secret/PHI logging.
- [ ] Pastikan migration dan constraint sesuai domain invariant.
- [ ] Pastikan semua dependency di-wire di composition root dan route dilindungi middleware/permission yang benar.
- [ ] Jalankan `gofmt`, test relevan/`go test ./...`, dan `go vet ./...` sejauh environment memungkinkan.
- [ ] Periksa race-sensitive code dan jalankan test concurrency/race yang relevan bila layak.
- [ ] Dokumentasikan asumsi, trade-off, migration/deployment order, dan test yang tidak dapat dijalankan.
- [ ] Jangan melakukan refactor di luar scope tanpa kebutuhan untuk correctness atau boundary yang sedang dikerjakan.

## 20. Aturan keputusan terakhir

Jika agent ragu antara jalur sederhana dan kompleks:

1. Jangan menilai dari jumlah endpoint atau tabel saja; nilai dari invariant, transaction, lifecycle, authorization, dan side effect.
2. Pertahankan CRUD sederhana tetap sederhana.
3. Untuk data klinis/finansial atau workflow lintas modul, pilih boundary yang lebih eksplisit dan aman.
4. Jika keputusan mengubah ownership modul, schema inti, Encounter lifecycle, atau jaminan transaksi, hentikan coding spekulatif dan minta keputusan arsitektural dengan opsi serta trade-off yang konkret.

Tujuan akhir bukan “Clean Architecture paling murni”, melainkan sistem yang business intent-nya mudah ditemukan, dependency-nya terkendali, transaksi medisnya konsisten, dan perubahan berikutnya tetap aman.
