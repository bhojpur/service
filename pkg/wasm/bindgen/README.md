# Bhojpur Service - WebAssemby BindGen

Let the `WebAssembly`'s exported function support more data types for its parameters and return values.

## Simple Example

Export Rust things to host program.

```rust
use wasmedge_bindgen::*;
use wasmedge_bindgen_macro::*;

// Export a `say` function from Rust, that returns a hello message.
#[wasmedge_bindgen]
pub fn say(s: String) -> String {
	let r = String::from("hello ");
	return r + s.as_str();
}
```

Use exported Rust things from WasmEdge!

```go
import (
	"github.com/bhojpur/service/pkg/wasm/wasmedge"
	bindgen "github.com/bhojpur/service/pkg/wasm/bindgen"
)

func main() {
	.
	.
	.

	// Instantiate the bindgen and vm
	bg := bindgen.Instantiate(vm)

		/// say: string -> string
	res, _ := bg.Execute("say", "wasmedge-bindgen")
	fmt.Println("Run bindgen -- say:", res[0].(string))
}
```
