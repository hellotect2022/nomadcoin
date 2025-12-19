**Go 언어 기반 프로젝트 학습 항목**

* **1. 파일 구조 및 의존성 관리**
    * `go.mod` : `package.json` 같은거, `requirements.txt`

* **2. Blockchain & Block 의 개념**
    * 1 blockchain & block 의 개념

* **3. [GO] Server Side Rendering Web Service**
    * - template 라이브러리 , handler, router
    * - json marshal, encoder
    * - TextMarshaler 인터페이스로 MarshalText 구현
    * - multiplexer (Mux, gorilla mux)
    * - error.New 로 새 에러메세지 만들기
    * - Middlewares 설정으로 application/json 미들웨어만들기

* **4. CLI Tool 만들기**
    * - flag 라이브러리 사용

* **5. [Go] BoltDB**
    * - blockchain 의 persistence 속성부여
    * - key-value DB 사용
    * - window 에서 boltbrowser 설치 방법 
    go install github.com/br0xen/boltbrowser@latest      
    boltbrowser blockChain.db
    * - boltweb 설치 및 실행방법
    go install github.com/evnix/boltdbweb@latest
    boltdbweb --db-name=blockChain.db --port=8888
    boltdbweb --db-name=<DBfilename>[required] --port=<port>[optional] --static-path=<static-path>[optional]

* **6. [POW & Mining] Proof of Work & Mining**
작업증명 : n 개의 0으로 시작하는 hash 를 찾는다
nonce : 사용자가 block 에서 변경가능한 유일한 값 
작업증명의 난이도를 알려주는 간단한 예시코드 
``` go
difficulty := 6
	target := strings.Repeat("0", difficulty)
	nonce := 1

	for {
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte("hello"+fmt.Sprint(nonce))))
		fmt.Println(hash)

		if strings.HasPrefix(hash, target) {
			fmt.Println("찾았습니당!!!", nonce)
			return
		} else {
			nonce++
		}
	}
```

채굴난이도 변경
비트코인 : 매 2016개 블록마다 얼마나 빠르게 블록이 생성되는지 계산(2주, 1개당 10분)

7 거래 
 - 비트코인 
- 무엇이 입력값을 생성하느가 ? Coinbase Transaction 
 Tx
  TxIn[$5(nico)]            : 거래하기 전에 너계좌에 있는 돈
  TxOut[$0(nico), $5(you)]  : 거래가 끝났을때 각각의 사람들이 가지고 있는돈 

exchange $5
  Tx $5
  TxIn[$10(nico)]
  TxOut[$5(you) $5(nico)]

  blocks : (느림) 
  mempool : (굉장히 빠름) 확정되지 않은 거래내역들을 기록해 두는 곳 
             확정되는 시점은 채굴자가 거래내역을 블록에 삽입하는 순간 (거래수수료 떼감, 채굴자는 다른 거래내역보다 먼저 처리해줌)

