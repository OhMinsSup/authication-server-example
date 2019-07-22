## lafu-server

[0] 서버 환경 설정

[] 회원가입 api 생성
    - accessToken, refreshToken 생성하는 함수
    - authToken 모델 생성
    - 리펙토링

## GO Modules 커맨드 정리

- go mod init [module-name]

  모듈을 생성합니다. 커맨드에서 [module-name]을 생략했다면, // import "[import-path]"를 추가하여야 합니다.

- go get [module-path]@[module-query]

  버전을 지정해 모듈을 추가합니다. [module-query]에 대해서는 공식 문서를 참고하시면 될 것 같습니다.

- go mod tidy [-v]

  go.mod 파일과 소스코드를 비교하여, import 되지 않은 의존성은 제거하고, import 되었지만 의존성 리스트에 추가되지 않은 모듈은 추가합니다. -v 플래그를 통해 더 자세한 정보를 확인할 수 있습니다.

- go mod vendor [-v]

  vendor/ 디렉터리를 생성합니다. -v 플래그를 통해 더 자세한 정보를 확인할 수 있습니다.

- go mod verify

  로컬에 설치된 모듈의 해시 값과 go.sum을 비교하여 모듈의 유효성을 검증합니다.
