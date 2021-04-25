start = path


path "path"
  = primary:primary ws "/" ws path:path {
    return {"and": [primary,path] }
  }
  / primary:primary ws "|" ws path:path {
    return {"or": [primary,path] }
  }
  / primary:primary {
    return primary;
  }

primary "primary"
  = iri:iri {
    return iri;
  }
  / "(" ws p:path ws ")" {
    return p;
  }


ws "whitespace"
  = [" "\n\t] * {
  return ""
}

iri "iri"
  = ns:[a-zA-Z0-9\-_]+ "." prop:[a-zA-Z0-9\-_]+ ws mod:["^","*"]? {
   var iri = {"iri": ns.join("") + ":" + prop.join(""), "inverse": false, "transitive": false};
   if (mod && mod === "^") {
     iri.inverse = true;
   } else if (mod && mod === "*") {
     iri.transitive = true;
   }
   return iri;
 }