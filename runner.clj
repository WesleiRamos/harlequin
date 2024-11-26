(ns runner
  (:require [joker.os :as os]))

(def project-properties (atom {}))

(defn defproject
  [& properties]
  (doseq [[property value] (partition 2 properties)]
    (cond
      (keyword? property)
      (swap! project-properties assoc property value)
      :else (throw (ex-info (str "Invalid property: " property) {})))))

(load-file "{{ PROJECT_FILE }}")

(let [source-path (:source-path @project-properties)
      main-ns (:main @project-properties)]
  (ns-sources {main-ns {:url source-path}})
  (require main-ns)
  (let [main-function (ns-resolve main-ns '-main)]
    (if main-function
      (apply main-function (os/args))
      (println "-main function not found" main-ns))))
